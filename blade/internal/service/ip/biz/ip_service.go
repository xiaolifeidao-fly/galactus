package biz

import (
	"errors"
	"galactus/blade/internal/service/ip"
	"galactus/blade/internal/service/ip/dto"
	"math/rand"
	"sync"
	"time"
)

var (
	ipInstance *IpService
	ipOnce     sync.Once
)

type IpService struct {
	ipEntities  []*dto.ProxyIpDTO
	ipNum       int
	mu          sync.Mutex
	baseService *ip.IpService
}

func GetDefaultIpService() *IpService {
	ipOnce.Do(func() {
		ipInstance = &IpService{
			baseService: ip.NewIpService(),
		}
		if err := ipInstance.InitIp(); err != nil {
			panic("Failed to initialize IpService: " + err.Error())
		}
	})
	return ipInstance
}

func (s *IpService) InitIp() error {
	proxyIps, err := s.baseService.GetProxyIpsByType("FINISH_QUERY")
	if err != nil || len(proxyIps) == 0 {
		return errors.New("not found ip config")
	}

	s.ipEntities = proxyIps
	s.ipNum = len(s.ipEntities)
	return nil
}

func (s *IpService) GetIp() (*dto.ProxyIpDTO, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ipDTO, err := s.getIpDTO()
	if err != nil {
		return nil, err
	}
	return ipDTO, nil
}

func (s *IpService) getIpDTO() (*dto.ProxyIpDTO, error) {
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

func (s *IpService) updateIp(ipDTO *dto.ProxyIpDTO, newIp map[string]interface{}) {
	ipDTO.Ip = newIp["ip"].(string) + ":" + newIp["port"].(string)
	ipDTO.ExpireTime = newIp["expireTime"].(time.Time)
	// Save updated IP back to the repository
	_, err := s.baseService.SaveOrUpdateProxyIp(ipDTO)
	if err != nil {
		// 这里可以添加日志记录
		return
	}
}

func (s *IpService) getRandomIndex() int {
	ipSize := len(s.ipEntities)
	if ipSize == 0 {
		return -1
	}
	return rand.Intn(ipSize)
}
