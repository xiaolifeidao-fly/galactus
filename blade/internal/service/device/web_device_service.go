package device

import (
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/device/repository"
	"galactus/common/middleware/db"
	"log"
)

type WebDeviceService struct {
	webDeviceRepository *repository.WebDeviceRepository
}

func NewWebDeviceService() *WebDeviceService {
	return &WebDeviceService{
		webDeviceRepository: db.GetRepository[repository.WebDeviceRepository](),
	}
}

// Save 保存设备信息
func (s *WebDeviceService) Save(webDeviceDTO *dto.WebDeviceDTO) (*dto.WebDeviceDTO, error) {
	device := db.ToPO[repository.WebDevice](webDeviceDTO)
	savedEntity, err := s.webDeviceRepository.SaveOrUpdate(device)
	if err != nil {
		return nil, err
	}
	return db.ToDTO[dto.WebDeviceDTO](savedEntity), nil
}

// GetById 根据ID获取设备
func (s *WebDeviceService) GetById(id uint) (*dto.WebDeviceDTO, error) {
	device, err := s.webDeviceRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return db.ToDTO[dto.WebDeviceDTO](device), nil
}

// GetActiveByRangeId 获取指定范围内的活跃设备
func (s *WebDeviceService) GetActiveByRangeId(startIndex, endIndex int64) ([]*dto.WebDeviceDTO, error) {
	devices, err := s.webDeviceRepository.GetActiveByRangeId(startIndex, endIndex)
	if err != nil {
		return nil, err
	}
	return db.ToDTOs[dto.WebDeviceDTO](devices), nil
}

// MinIdByStartIndex 获取起始索引之后的最小ID
func (s *WebDeviceService) MinIdByStartIndex(startIndex int64) (int64, error) {
	return s.webDeviceRepository.MinIdByStartIndex(startIndex)
}

// GetByUdIdAndOpenUdIdAndSerial 根据设备唯一标识获取设备
func (s *WebDeviceService) GetByUdIdAndOpenUdIdAndSerial(udId, openUdId, serial string) (*dto.WebDeviceDTO, error) {
	device, err := s.webDeviceRepository.GetByUdIdAndOpenUdIdAndSerial(udId, openUdId, serial)
	if err != nil {
		return nil, err
	}
	return db.ToDTO[dto.WebDeviceDTO](device), nil
}

// GetWebDeviceList 获取所有Web设备列表
func (s *WebDeviceService) GetWebDeviceList() ([]*dto.WebDeviceDTO, error) {
	devices, err := s.webDeviceRepository.FindAll()
	if err != nil {
		log.Printf("get web device list failed: %v", err)
		return nil, err
	}
	return db.ToDTOs[dto.WebDeviceDTO](devices), nil
}
