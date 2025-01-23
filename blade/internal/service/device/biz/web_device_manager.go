package biz

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"galactus/blade/internal/consts"
	"galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/dictionary"
	"galactus/blade/internal/service/ip/biz"
	"galactus/common/middleware/vipper"
)

var (
	defaultWebDeviceManager *WebDeviceManager
	webDeviceManagerOnce    sync.Once
)

type WebDeviceManager struct {
	mu sync.Mutex
	// 使用 map 存储不同场景的设备列表
	devicesMap         map[consts.Scene]*SceneDevices
	webDeviceSvc       *device.WebDeviceService
	dictionarySvc      dictionary.DictionaryService
	webDeviceExpireKey string
	ipManager          *biz.IpManager
	deviceIpMap        map[string][]*dto.WebDeviceDTO
}

// SceneDevices 每个场景的设备状态
type SceneDevices struct {
	devices      []*dto.WebDeviceDTO // 设备列表
	currentIndex int                 // 当前获取设备的索引
	minId        int                 // 设备最小ID
	maxId        int                 // 设备最大ID
	idRange      int                 // 设备范围
	currentLoop  int                 // 当前循环次数
	maxUse       int                 // 设备最大使用次数
}

// getWebDeviceRange 获取设备范围配置
func (m *WebDeviceManager) getWebDeviceRange(scene consts.Scene) (int64, error) {
	webDeviceRange, err := m.dictionarySvc.GetByCode(scene.GetDeviceIdRange())
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
			devicesMap:         make(map[consts.Scene]*SceneDevices),
			webDeviceSvc:       device.NewWebDeviceService(),
			dictionarySvc:      dictionary.NewDictionaryService(),
			webDeviceExpireKey: "WEB_DEVICE_EXPIRE",
			ipManager:          biz.GetDefaultIpManager(),
			deviceIpMap:        make(map[string][]*dto.WebDeviceDTO),
		}

		// 初始化所有场景
		for _, scene := range []consts.Scene{
			consts.SceneCollectDevice,
			consts.SceneAuditLike,
			consts.SceneCurrentValue,
			consts.SceneAuditFollow,
		} {
			// 获取场景的最大使用次数
			maxUseConfig, err := manager.dictionarySvc.GetByCode(scene.GetDeviceMaxUse())
			if err != nil {
				initErr = fmt.Errorf("get device max use config error for scene %v: %v", scene, err)
				return
			}

			maxUse, err := strconv.Atoi(maxUseConfig.Value)
			if err != nil {
				initErr = fmt.Errorf("parse max use config error for scene %v: %v", scene, err)
				return
			}

			minIdConfig, err := manager.dictionarySvc.GetByCode(scene.GetDeviceMinId())
			if err != nil {
				initErr = fmt.Errorf("get device min id config error for scene %v: %v", scene, err)
				return
			}

			minId, err := strconv.Atoi(minIdConfig.Value)
			if err != nil {
				initErr = fmt.Errorf("parse min id config error for scene %v: %v", scene, err)
				return
			}

			maxIdConfig, err := manager.dictionarySvc.GetByCode(scene.GetDeviceMaxId())
			if err != nil {
				initErr = fmt.Errorf("get device max id config error for scene %v: %v", scene, err)
				return
			}

			maxId, err := strconv.Atoi(maxIdConfig.Value)
			if err != nil {
				initErr = fmt.Errorf("parse max id config error for scene %v: %v", scene, err)
				return
			}

			idRangeConfig, err := manager.dictionarySvc.GetByCode(scene.GetDeviceIdRange())
			if err != nil {
				initErr = fmt.Errorf("get device id range config error for scene %v: %v", scene, err)
				return
			}

			idRange, err := strconv.Atoi(idRangeConfig.Value)
			if err != nil {
				initErr = fmt.Errorf("parse id range config error for scene %v: %v", scene, err)
				return
			}

			manager.devicesMap[scene] = &SceneDevices{
				devices:      make([]*dto.WebDeviceDTO, 0),
				currentIndex: 0,
				currentLoop:  0,
				maxUse:       maxUse,
				minId:        minId,
				maxId:        maxId,
				idRange:      idRange,
			}
		}

		defaultWebDeviceManager = manager
		defaultWebDeviceManager.ipManager.RegisterObserver(defaultWebDeviceManager)

		// 初始化所有场景的设备池
		for scene := range manager.devicesMap {
			if err := defaultWebDeviceManager.InitWebDevicePool(scene); err != nil {
				log.Printf("Failed to initialize WebDeviceManager for scene %v: %v", scene, err)
			}
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

// InitWebDevicePool 初始化指定场景的设备池
func (m *WebDeviceManager) InitWebDevicePool(scene consts.Scene) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sceneDevices, exists := m.devicesMap[scene]
	if !exists {
		return fmt.Errorf("scene %v not initialized", scene)
	}

	// 获取当前索引
	currentIndexConfig, err := m.dictionarySvc.GetByCode(scene.GetDeviceCurrentIndex())
	if err != nil {
		return fmt.Errorf("get current index error: %v", err)
	}
	currentIndex, err := strconv.Atoi(currentIndexConfig.Value)
	if err != nil {
		return fmt.Errorf("parse current index error: %v", err)
	}

	// 确保当前索引在有效范围内
	if currentIndex < sceneDevices.minId {
		currentIndex = sceneDevices.minId
	}
	if currentIndex > sceneDevices.maxId {
		currentIndex = sceneDevices.minId // 如果超出最大ID，重置为最小ID
	}

	// 获取设备列表，使用场景配置的idRange作为pageSize
	devices, err := m.webDeviceSvc.GetActiveByStartAndLimitWithRange(int64(currentIndex), int64(sceneDevices.idRange), int64(sceneDevices.maxId))
	if err != nil {
		return fmt.Errorf("get web devices error: %v", err)
	}

	if len(devices) == 0 {
		// 如果没有获取到设备，从最小ID重新开始
		devices, err = m.webDeviceSvc.GetActiveByStartAndLimitWithRange(int64(sceneDevices.minId), int64(sceneDevices.idRange), int64(sceneDevices.maxId))
		if err != nil {
			return fmt.Errorf("get web devices error: %v", err)
		}
	}

	if len(devices) > 0 {
		nextIndex := int64(devices[len(devices)-1].Id)
		if nextIndex >= int64(sceneDevices.maxId) {
			nextIndex = int64(sceneDevices.minId) // 如果达到最大ID，重置为最小ID
		}
		err = m.updateCurrentIndex(scene, nextIndex)
		if err != nil {
			return fmt.Errorf("update current index error: %v", err)
		}
		sceneDevices.devices = devices
		sceneDevices.currentIndex = 0
		sceneDevices.currentLoop = 0
	}

	return nil
}

// GetWebDevice 获取指定场景的可用设备
func (m *WebDeviceManager) GetWebDevice(scene consts.Scene) (*dto.WebDeviceDTO, error) {
	m.mu.Lock()
	sceneDevices, exists := m.devicesMap[scene]
	if !exists {
		m.mu.Unlock()
		return nil, fmt.Errorf("scene %v not initialized", scene)
	}

	// 如果没有设备，初始化设备池
	if len(sceneDevices.devices) == 0 {
		m.mu.Unlock()
		if err := m.InitWebDevicePool(scene); err != nil {
			return nil, fmt.Errorf("failed to initialize device pool: %v", err)
		}
		m.mu.Lock()
		if len(sceneDevices.devices) == 0 {
			m.mu.Unlock()
			return nil, fmt.Errorf("no available web device for scene %v", scene)
		}
	}

	// 如果已经遍历完所有设备
	if sceneDevices.currentIndex >= len(sceneDevices.devices) {
		sceneDevices.currentIndex = 0
		sceneDevices.currentLoop++

		// 如果已经达到最大循环次数，重新获取设备列表
		if sceneDevices.currentLoop >= int(sceneDevices.maxUse) {
			m.mu.Unlock()
			if err := m.InitWebDevicePool(scene); err != nil {
				return nil, fmt.Errorf("failed to reinitialize device pool: %v", err)
			}
			m.mu.Lock()
			if len(sceneDevices.devices) == 0 {
				m.mu.Unlock()
				return nil, fmt.Errorf("no available web device for scene %v", scene)
			}
		}
	}

	dev := sceneDevices.devices[sceneDevices.currentIndex]
	sceneDevices.currentIndex++
	m.mu.Unlock()

	m.fillWebDeviceDefaultValues(dev, scene)
	return dev, nil
}

// fillWebDeviceDefaultValues 填充设备默认值
func (m *WebDeviceManager) fillWebDeviceDefaultValues(dev *dto.WebDeviceDTO, scene consts.Scene) {
	if dev.Platform == "" {
		dev.Platform = "android"
	}
	if dev.UserAgent == "" {
		dev.UserAgent = ""
	}

	if vipper.GetBool("proxy.enable") {
		if dev.ProxyIp == "" {
			if ipDTO, err := m.ipManager.GetIp(scene); err == nil {
				dev.ProxyIp = ipDTO.Ip
				dev.ExpireTime = ipDTO.ExpireTime
				m.mu.Lock()
				m.deviceIpMap[ipDTO.Ip] = append(m.deviceIpMap[ipDTO.Ip], dev)
				m.mu.Unlock()
			}
		}
	}
}

// updateCurrentIndex 更新当前索引
func (m *WebDeviceManager) updateCurrentIndex(scene consts.Scene, index int64) error {
	config, err := m.dictionarySvc.GetByCode(scene.GetDeviceCurrentIndex())
	if err != nil {
		return err
	}

	config.Value = strconv.FormatInt(index, 10)
	return m.dictionarySvc.SaveOrUpdate(config)
}
