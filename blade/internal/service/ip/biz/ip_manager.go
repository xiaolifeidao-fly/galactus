package biz

import (
	"errors"
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
	proxyIps, err := s.baseService.GetProxyIpsByType("FINISH_QUERY")
	if err != nil || len(proxyIps) == 0 {
		log.Printf("not found ip config")
		return nil
	}

	s.ipEntities = proxyIps
	s.ipNum = len(s.ipEntities)
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
	randomIndex := s.getRandomIndex()
	if randomIndex == -1 {
		return nil, errors.New("no available IP")
	}
	ipDTO := s.ipEntities[randomIndex]
	now := time.Now()

	if ipDTO.ExpireTime.Before(now) {
		// Update IP logic here
		newIp, err := ip.GetDefaultZDYHttpProxyService().GetUserIpByProxyType(s.ipNum)
		if err != nil {
			return nil, err
		}
		s.updateIp(ipDTO, newIp)
	}
	return ipDTO, nil
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

func (s *IpManager) updateIp(ipDTO *dto.ProxyIpDTO, newIp map[string]interface{}) {
	oldIp := ipDTO.Ip
	ipDTO.Ip = newIp["ip"].(string) + ":" + newIp["port"].(string)
	ipDTO.ExpireTime = newIp["expireTime"].(time.Time)
	// Save updated IP back to the repository
	_, err := s.baseService.SaveOrUpdateProxyIp(ipDTO)
	if err != nil {
		log.Printf("update ip failed: %v", err)
		return
	}

	// 通知观察者IP已更新
	s.notifyObservers(oldIp, ipDTO.Ip)
}

func (s *IpManager) getRandomIndex() int {
	ipSize := len(s.ipEntities)
	if ipSize == 0 {
		return -1
	}
	return rand.Intn(ipSize)
}
