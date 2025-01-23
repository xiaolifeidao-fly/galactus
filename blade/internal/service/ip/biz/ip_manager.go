package biz

import (
	"errors"
	"fmt"
	"galactus/blade/internal/consts"
	"galactus/blade/internal/service/ip"
	"galactus/blade/internal/service/ip/dto"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	defaultIpInstance *IpManager
	ipManagerOnce     sync.Once
	maxRetries        = 3 // 最大重试次数
)

// IpObserver 定义IP更新的观察者接口
type IpObserver interface {
	OnIpUpdate(oldIp, newIp string)
}

type IpManager struct {
	// 使用 map 存储不同场景的 IP 列表
	ipEntitiesMap map[consts.Scene][]*dto.ProxyIpDTO
	ipNum         int
	mu            sync.Mutex
	baseService   *ip.IpService
	observers     []IpObserver
}

// GetDefaultIpManager 只负责获取实例，不负责初始化数据
func GetDefaultIpManager() *IpManager {
	ipManagerOnce.Do(func() {
		if defaultIpInstance == nil {
			defaultIpInstance = &IpManager{
				baseService:   ip.NewIpService(),
				observers:     make([]IpObserver, 0),
				ipNum:         10,
				ipEntitiesMap: make(map[consts.Scene][]*dto.ProxyIpDTO),
			}
		}
	})
	return defaultIpInstance
}

// InitIpManager 显式初始化方法，包含数据加载
func InitIpManager() error {
	manager := GetDefaultIpManager()
	return manager.InitIp()
}

func (s *IpManager) InitIp() error {
	proxyIps, err := s.baseService.GetAllProxyIps()
	if err != nil || len(proxyIps) == 0 {
		log.Printf("not found ip config")
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 初始化map
	s.ipEntitiesMap = make(map[consts.Scene][]*dto.ProxyIpDTO)

	// 根据Type将IP分配到对应的场景
	for _, ip := range proxyIps {
		// 遍历所有场景，找到匹配的场景
		for _, scene := range []consts.Scene{
			consts.SceneCollectDevice,
			consts.SceneAuditLike,
			consts.SceneCurrentValue,
			consts.SceneAuditFollow,
		} {
			if ip.Type == scene.GetSceneName() {
				s.ipEntitiesMap[scene] = append(s.ipEntitiesMap[scene], ip)
				break
			}
		}
	}
	return nil
}

func (s *IpManager) GetIp(scene consts.Scene) (*dto.ProxyIpDTO, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ipDTO, err := s.getIpDTO(scene)
	if err != nil {
		return nil, err
	}
	return ipDTO, nil
}

func (s *IpManager) getIpDTO(scene consts.Scene) (*dto.ProxyIpDTO, error) {
	retryCount := 0
	for {
		if retryCount >= maxRetries {
			return nil, errors.New("exceeded maximum retry attempts to get valid IP")
		}

		// 获取指定场景的 IP 列表
		ipEntities, exists := s.ipEntitiesMap[scene]
		if !exists {
			s.ipEntitiesMap[scene] = make([]*dto.ProxyIpDTO, 0)
			ipEntities = s.ipEntitiesMap[scene]
		}

		if len(ipEntities) == 0 {
			// 获取新的 IP 列表
			proxyIPs, err := ip.GetDefaultZDYHttpProxyService().GetUserIpByProxyType(scene, s.ipNum)
			if err != nil {
				retryCount++
				log.Printf("failed to get IPs from ZDY service for scene %v, retry %d/%d: %v",
					scene, retryCount, maxRetries, err)
				continue
			}

			// 保存新的 IP 列表
			newIpEntities := make([]*dto.ProxyIpDTO, 0, len(proxyIPs))
			for _, proxyIP := range proxyIPs {
				ipDTO := &dto.ProxyIpDTO{
					Ip:         proxyIP.IP + ":" + fmt.Sprint(proxyIP.Port),
					ExpireTime: time.Now().Add(time.Duration(proxyIP.Timeout) * time.Second),
					Type:       scene.GetSceneName(),
				}
				savedDTO, err := s.baseService.SaveOrUpdateProxyIp(ipDTO)
				if err != nil {
					log.Printf("failed to save proxy IP: %v", err)
					continue
				}
				newIpEntities = append(newIpEntities, savedDTO)
			}

			if len(newIpEntities) == 0 {
				retryCount++
				log.Printf("no valid IPs saved for scene %v, retry %d/%d", scene, retryCount, maxRetries)
				continue
			}

			s.ipEntitiesMap[scene] = newIpEntities
			ipEntities = newIpEntities
		}

		// 从当前场景的 IP 列表中随机选择一个
		randomIndex := rand.Intn(len(ipEntities))
		ipDTO := ipEntities[randomIndex]
		now := time.Now()

		if ipDTO.ExpireTime.Before(now) {
			// 删除过期的 IP
			if err := s.baseService.DeleteProxyIp(int64(ipDTO.Id)); err != nil {
				log.Printf("failed to delete expired IP from database: %v", err)
			}
			// 从内存中移除过期的 IP
			s.ipEntitiesMap[scene] = append(ipEntities[:randomIndex], ipEntities[randomIndex+1:]...)
			continue
		}

		return ipDTO, nil
	}
}

// RegisterObserver 注册观察者
func (s *IpManager) RegisterObserver(observer IpObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

// notifyObservers 通知所有观察者
func (s *IpManager) notifyObservers(oldIp, newIp string) {
	for _, observer := range s.observers {
		go observer.OnIpUpdate(oldIp, newIp)
	}
}
