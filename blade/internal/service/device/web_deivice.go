package device

import (
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/device/repository"
	"galactus/common/middleware/db"
)

var webDeviceRepository *repository.WebDeviceRepository = db.GetRepository[repository.WebDeviceRepository]()

func SaveWebDevice(deviceDTO *dto.WebDevice) (*dto.WebDevice, error) {
	device := db.ToPO[repository.WebDevice](deviceDTO)
	device, err := webDeviceRepository.SaveOrUpdate(device)
	if err != nil {
		return nil, err
	}
	return db.ToDTO[dto.WebDevice](device), nil
}

func ListWebDevice() ([]*dto.WebDevice, error) {
	devices, err := webDeviceRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return db.ToDTOs[dto.WebDevice](devices), nil
}
