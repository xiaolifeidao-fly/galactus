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

type IpManager struct {
	ipEntities  []*dto.ProxyIpDTO
	ipNum       int
	mu          sync.Mutex
	baseService *ip.IpService
}

func InitDefaultIpManager() {
	ipManagerOnce.Do(func() {
		defaultIpInstance = &IpManager{
			baseService: ip.NewIpService(),
		}
		if err := defaultIpInstance.InitIp(); err != nil {
			log.Printf("Failed to initialize IpManager: %v", err)
		}
	})
}

func GetDefaultIpManager() *IpManager {
	if defaultIpInstance == nil {
		panic("IpManager is not initialized. Call InitDefaultIpManager first.")
	}
	return defaultIpInstance
}

func (s *IpManager) InitIp() error {
	proxyIps, err := s.baseService.GetProxyIpsByType("FINISH_QUERY")
	if err != nil || len(proxyIps) == 0 {
		return errors.New("not found ip config")
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

func (s *IpManager) updateIp(ipDTO *dto.ProxyIpDTO, newIp map[string]interface{}) {
	ipDTO.Ip = newIp["ip"].(string) + ":" + newIp["port"].(string)
	ipDTO.ExpireTime = newIp["expireTime"].(time.Time)
	// Save updated IP back to the repository
	_, err := s.baseService.SaveOrUpdateProxyIp(ipDTO)
	if err != nil {
		// 这里可以添加日志记录
		log.Printf("update ip failed: %v", err)
		return
	}
}

func (s *IpManager) getRandomIndex() int {
	ipSize := len(s.ipEntities)
	if ipSize == 0 {
		return -1
	}
	return rand.Intn(ipSize)
}
