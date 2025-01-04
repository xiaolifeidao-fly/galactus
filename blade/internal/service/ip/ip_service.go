package ip

import (
	"errors"
	"galactus/blade/internal/service/ip/dto"
	"galactus/blade/internal/service/ip/repository"
	"galactus/common/middleware/db"
	"math/rand"
	"sync"
	"time"
)

var (
	ipInstance   *IpService
	ipOnce       sync.Once
	ipRepository *repository.ProxyIpRepository
	repoInitOnce sync.Once
)

type IpService struct {
	ipEntities []*dto.ProxyIpDTO
	ipNum      int
	mu         sync.Mutex
}

func InitDefaultIpService() (*IpService, error) {
	var err error
	ipOnce.Do(func() {
		repoInitOnce.Do(func() {
			ipRepository = db.GetRepository[repository.ProxyIpRepository]()
		})
		ipInstance = &IpService{}
		err = ipInstance.InitIp()
	})
	return ipInstance, err
}

func GetDefaultIpService() *IpService {
	if ipInstance == nil {
		panic("IpService is not initialized. Call InitDefaultIpService first.")
	}
	return ipInstance
}

func (s *IpService) InitIp() error {
	proxyIps, err := ipRepository.GetByType("FINISH_QUERY")
	if err != nil || len(proxyIps) == 0 {
		return errors.New("not found ip config")
	}

	// Convert entities to DTOs
	s.ipEntities = db.ToDTOs[dto.ProxyIpDTO](proxyIps)
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
		newIp, err := GetDefaultZDYHttpProxyService().GetUserIpByProxyType(s.ipNum)
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
	ipRepository.SaveOrUpdate(db.ToPO[repository.ProxyIp](ipDTO))
}

func (s *IpService) getRandomIndex() int {
	ipSize := len(s.ipEntities)
	if ipSize == 0 {
		return -1
	}
	return rand.Intn(ipSize)
}
