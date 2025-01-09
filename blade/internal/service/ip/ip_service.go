package ip

import (
	"galactus/blade/internal/service/ip/dto"
	"galactus/blade/internal/service/ip/repository"
	"galactus/common/middleware/db"
)

type IpService struct {
	ipRepository *repository.ProxyIpRepository
}

func NewIpService() *IpService {
	return &IpService{
		ipRepository: db.GetRepository[repository.ProxyIpRepository](),
	}
}

// GetProxyIpsByType 根据类型获取代理IP
func (s *IpService) GetProxyIpsByType(proxyType string) ([]*dto.ProxyIpDTO, error) {
	proxyIps, err := s.ipRepository.GetByType(proxyType)
	if err != nil {
		return nil, err
	}
	return db.ToDTOs[dto.ProxyIpDTO](proxyIps), nil
}

// GetAllProxyIps 获取所有代理IP
func (s *IpService) GetAllProxyIps() ([]*dto.ProxyIpDTO, error) {
	proxyIps, err := s.ipRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return db.ToDTOs[dto.ProxyIpDTO](proxyIps), nil
}

// SaveOrUpdateProxyIp 保存或更新代理IP
func (s *IpService) SaveOrUpdateProxyIp(proxyIpDTO *dto.ProxyIpDTO) (*dto.ProxyIpDTO, error) {
	proxyIp := db.ToPO[repository.ProxyIp](proxyIpDTO)
	savedEntity, err := s.ipRepository.SaveOrUpdate(proxyIp)
	if err != nil {
		return nil, err
	}
	return db.ToDTO[dto.ProxyIpDTO](savedEntity), nil
}

// DeleteProxyIp 删除代理IP
func (s *IpService) DeleteProxyIp(id int64) error {
	return s.ipRepository.Delete(uint(id))
}
