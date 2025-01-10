package biz

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/dictionary"
	"galactus/blade/internal/service/dictionary/constants"
	"galactus/blade/internal/service/ip/biz"
	"galactus/common/middleware/vipper"
)

var (
	defaultWebDeviceManager *WebDeviceManager
	webDeviceManagerOnce    sync.Once
)

type WebDeviceManager struct {
	mu                 sync.Mutex
	devices            []*dto.WebDeviceDTO // 设备列表
	currentIndex       int                 // 当前获取设备的索引
	currentLoop        int                 // 当前循环次数
	maxUse             int64               // 设备最大使用次数
	webDeviceSvc       *device.WebDeviceService
	dictionarySvc      dictionary.DictionaryService
	webDeviceExpireKey string
	ipManager          *biz.IpManager
	deviceIpMap        map[string][]*dto.WebDeviceDTO
}

// getWebDeviceRange 获取设备范围配置
func (m *WebDeviceManager) getWebDeviceRange() (int64, error) {
	webDeviceRange, err := m.dictionarySvc.GetByCode(constants.WEB_DEVICE_ID_RANGE.Code)
	if err != nil {
		return 0, fmt.Errorf("get web device range config error: %v", err)
	}
	return strconv.ParseInt(webDeviceRange.Value, 10, 64)
}

// InitDefaultWebDeviceManager 初始化默认的WebDeviceManager实例
func InitDefaultWebDeviceManager() error {
	var initErr error
	webDeviceManagerOnce.Do(func() {
		manager := &WebDeviceManager{
			devices:            make([]*dto.WebDeviceDTO, 0),
			currentIndex:       0,
			currentLoop:        0,
			webDeviceSvc:       device.NewWebDeviceService(),
			dictionarySvc:      dictionary.NewDictionaryService(),
			webDeviceExpireKey: "WEB_DEVICE_EXPIRE",
			ipManager:          biz.GetDefaultIpManager(),
			deviceIpMap:        make(map[string][]*dto.WebDeviceDTO),
		}

		// 获取设备最大使用次数配置
		maxUseConfig, err := manager.dictionarySvc.GetByCode(constants.WEB_DEVICE_MAX_USE.Code)
		if err != nil {
			initErr = fmt.Errorf("get web device max use config error: %v", err)
			return
		}
		maxUse, err := strconv.ParseInt(maxUseConfig.Value, 10, 64)
		if err != nil {
			initErr = fmt.Errorf("parse max use config error: %v", err)
			return
		}
		manager.maxUse = maxUse

		defaultWebDeviceManager = manager
		// 注册为IP更新的观察者
		defaultWebDeviceManager.ipManager.RegisterObserver(defaultWebDeviceManager)

		// 初始化设备池
		if err := defaultWebDeviceManager.InitWebDevicePool(); err != nil {
			log.Printf("Failed to initialize WebDeviceManager: %v", err)
		}
	})
	return initErr
}

// GetDefaultWebDeviceManager 获取默认的WebDeviceManager实例
func GetDefaultWebDeviceManager() *WebDeviceManager {
	if defaultWebDeviceManager == nil {
		panic("WebDeviceManager is not initialized. Call InitDefaultWebDeviceManager first.")
	}
	return defaultWebDeviceManager
}

// OnIpUpdate 实现IP更新的观察者接口
func (m *WebDeviceManager) OnIpUpdate(oldIp, newIp string) {
	// m.mu.Lock()
	// defer m.mu.Unlock()

	// // 查找使用旧IP的设备
	// if devices, ok := m.deviceIpMap[oldIp]; ok {
	// 	// 更新这些设备的IP
	// 	for _, dev := range devices {
	// 		dev.ProxyIp = newIp
	// 	}
	// 	// 更新映射关系
	// 	delete(m.deviceIpMap, oldIp)
	// 	m.deviceIpMap[newIp] = devices
	// }
}

// InitWebDevicePool 初始化设备池
func (m *WebDeviceManager) InitWebDevicePool() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取设备范围
	pageSize, err := m.getWebDeviceRange()
	if err != nil {
		return err
	}

	// 获取当前索引
	currentIndex, err := m.getCurrentIndex()
	if err != nil {
		return fmt.Errorf("get current index error: %v", err)
	}

	// 获取设备列表
	devices, err := m.webDeviceSvc.GetActiveByStartAndLimit(currentIndex, pageSize)
	if err != nil {
		return fmt.Errorf("get web devices error: %v", err)
	}

	if len(devices) == 0 {
		//循环重新获取
		devices, err = m.webDeviceSvc.GetActiveByStartAndLimit(0, pageSize)
		if err != nil {
			return fmt.Errorf("get web devices error: %v", err)
		}
	}

	if len(devices) > 0 {
		err = m.updateCurrentIndex(int64(devices[len(devices)-1].Id))
		if err != nil {
			return fmt.Errorf("update current index error: %v", err)
		}
		// 更新设备列表
		m.devices = devices
		m.currentIndex = 0
		m.currentLoop = 0 // 重置循环次数
	}

	return nil
}

// GetWebDevice 获取一个可用设备
func (m *WebDeviceManager) GetWebDevice() (*dto.WebDeviceDTO, error) {
	m.mu.Lock()
	// 如果没有设备，初始化设备池
	if len(m.devices) == 0 {
		m.mu.Unlock()
		if err := m.InitWebDevicePool(); err != nil {
			return nil, fmt.Errorf("failed to initialize device pool: %v", err)
		}
		m.mu.Lock()
		if len(m.devices) == 0 {
			m.mu.Unlock()
			return nil, fmt.Errorf("no available web device")
		}
	}

	// 如果已经遍历完所有设备
	if m.currentIndex >= len(m.devices) {
		m.currentIndex = 0 // 重置索引
		m.currentLoop++    // 增加循环次数

		// 如果已经达到最大循环次数，重新获取设备列表
		if m.currentLoop >= int(m.maxUse) {
			m.mu.Unlock()
			if err := m.InitWebDevicePool(); err != nil {
				return nil, fmt.Errorf("failed to reinitialize device pool: %v", err)
			}
			m.mu.Lock()
			if len(m.devices) == 0 {
				m.mu.Unlock()
				return nil, fmt.Errorf("no available web device")
			}
		}
	}

	// 获取当前设备
	dev := m.devices[m.currentIndex]
	m.currentIndex++
	m.mu.Unlock()

	// 设置默认值
	m.fillWebDeviceDefaultValues(dev)

	return dev, nil
}

// fillWebDeviceDefaultValues 填充设备默认值
func (m *WebDeviceManager) fillWebDeviceDefaultValues(dev *dto.WebDeviceDTO) {
	if dev.Platform == "" {
		dev.Platform = "android"
	}
	if dev.UserAgent == "" {
		dev.UserAgent = ""
	}

	if vipper.GetBool("proxy.enable") {
		// 设置IP
		if dev.ProxyIp == "" {
			if ipDTO, err := m.ipManager.GetIp(); err == nil {
				dev.ProxyIp = ipDTO.Ip
				// 记录IP和设备的关系
				m.mu.Lock()
				m.deviceIpMap[ipDTO.Ip] = append(m.deviceIpMap[ipDTO.Ip], dev)
				m.mu.Unlock()
			}
		}
	}

}

// getCurrentIndex 获取当前索引
func (m *WebDeviceManager) getCurrentIndex() (int64, error) {
	config, err := m.dictionarySvc.GetByCode(constants.WEB_DEVICE_CURRENT_INDEX.Code)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(config.Value, 10, 64)
}

// updateCurrentIndex 更新当前索引
func (m *WebDeviceManager) updateCurrentIndex(index int64) error {
	config, err := m.dictionarySvc.GetByCode(constants.WEB_DEVICE_CURRENT_INDEX.Code)
	if err != nil {
		return err
	}

	config.Value = strconv.FormatInt(index, 10)
	return m.dictionarySvc.SaveOrUpdate(config)
}
