package device

import (
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/device/repository"
	"galactus/common/middleware/db"
)

var webDeviceRepository *repository.WebDeviceRepository = db.GetRepository[repository.WebDeviceRepository]()

func SaveWebDevice(deviceDTO *dto.WebDeviceDTO) (*dto.WebDeviceDTO, error) {
	device := db.ToPO[repository.WebDevice](deviceDTO)
	device, err := webDeviceRepository.SaveOrUpdate(device)
	if err != nil {
		return nil, err
	}
	return db.ToDTO[dto.WebDeviceDTO](device), nil
}

func GetWebDeviceList() ([]*dto.WebDeviceDTO, error) {
	devices, err := webDeviceRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return db.ToDTOs[dto.WebDeviceDTO](devices), nil
}

func FindWebDeviceById(id uint) (*dto.WebDeviceDTO, error) {
	device, err := webDeviceRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return db.ToDTO[dto.WebDeviceDTO](device), nil
}
