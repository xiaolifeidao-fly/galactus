package biz

import (
	"errors"
	"fmt"
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
	ipEntities  []*dto.ProxyIpDTO
	ipNum       int
	mu          sync.Mutex
	baseService *ip.IpService
	observers   []IpObserver
}

// GetDefaultIpManager 只负责获取实例，不负责初始化数据
func GetDefaultIpManager() *IpManager {
	ipManagerOnce.Do(func() {
		if defaultIpInstance == nil {
			defaultIpInstance = &IpManager{
				baseService: ip.NewIpService(),
				observers:   make([]IpObserver, 0),
				ipNum:       10, //每次拿10个ip
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

	s.ipEntities = proxyIps
	return nil
}

func (s *IpManager) GetIp() (*dto.ProxyIpDTO, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ipDTO, err := s.getIpDTO()
	if err != nil {
		return nil, err
	}
	return ipDTO, nil
}

func (s *IpManager) getIpDTO() (*dto.ProxyIpDTO, error) {
	retryCount := 0
	for {
		if retryCount >= maxRetries {
			return nil, errors.New("exceeded maximum retry attempts to get valid IP")
		}

		if len(s.ipEntities) == 0 {
			// 如果没有可用IP，重新获取新的IP列表
			proxyIPs, err := ip.GetDefaultZDYHttpProxyService().GetUserIpByProxyType(s.ipNum)
			if err != nil {
				retryCount++
				log.Printf("failed to get IPs from ZDY service, retry %d/%d: %v", retryCount, maxRetries, err)
				continue
			}

			// 将新的IP列表转换为DTO并保存
			for _, proxyIP := range proxyIPs {
				ipDTO := &dto.ProxyIpDTO{
					Ip:         proxyIP.IP + ":" + fmt.Sprint(proxyIP.Port),
					ExpireTime: time.Now().Add(1 * time.Minute), // 设置固定5分钟失效时间
				}
				savedDTO, err := s.baseService.SaveOrUpdateProxyIp(ipDTO)
				if err != nil {
					log.Printf("failed to save proxy IP: %v", err)
					continue
				}
				s.ipEntities = append(s.ipEntities, savedDTO)
			}

			if len(s.ipEntities) == 0 {
				retryCount++
				log.Printf("no valid IPs saved, retry %d/%d", retryCount, maxRetries)
				continue
			}
		}

		randomIndex := s.getRandomIndex()
		if randomIndex == -1 {
			retryCount++
			continue
		}

		ipDTO := s.ipEntities[randomIndex]
		now := time.Now()

		if ipDTO.ExpireTime.Before(now) {
			// 从数据库中删除过期IP
			if err := s.baseService.DeleteProxyIp(int64(ipDTO.Id)); err != nil {
				log.Printf("failed to delete expired IP from database: %v", err)
			}
			// 从内存中移除过期的IP
			s.ipEntities = append(s.ipEntities[:randomIndex], s.ipEntities[randomIndex+1:]...)
			continue // 继续循环尝试下一个IP
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

func (s *IpManager) getRandomIndex() int {
	ipSize := len(s.ipEntities)
	if ipSize == 0 {
		return -1
	}
	return rand.Intn(ipSize)
}
