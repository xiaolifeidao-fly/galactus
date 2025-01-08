package biz

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/dictionary"
	"galactus/blade/internal/service/dictionary/constants"
	"galactus/blade/internal/service/ip/biz"
	"galactus/common/middleware/redis"
)

const (
	webDeviceExpireTime = 30 * 60 // 30分钟
	webDevicePoolKey    = "WEB_DEVICE_POOL"
	deviceTimeout       = 5 * time.Second // 固定5秒超时
)

var (
	defaultWebDeviceManager *WebDeviceManager
	webDeviceManagerOnce    sync.Once
)

type WebDeviceManager struct {
	mu                 sync.Mutex
	webDevicePool      chan *dto.WebDeviceDTO
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
	webDeviceManagerOnce.Do(func() {
		// 获取设备范围配置
		defaultWebDeviceManager = &WebDeviceManager{
			webDevicePool:      make(chan *dto.WebDeviceDTO, 1000), // 使用固定大小作为channel容量
			webDeviceSvc:       device.NewWebDeviceService(),
			dictionarySvc:      dictionary.NewDictionaryService(),
			webDeviceExpireKey: "WEB_DEVICE_EXPIRE",
			ipManager:          biz.GetDefaultIpManager(),
			deviceIpMap:        make(map[string][]*dto.WebDeviceDTO),
		}
		// 注册为IP更新的观察者
		defaultWebDeviceManager.ipManager.RegisterObserver(defaultWebDeviceManager)

		// 初始化设备池
		if err := defaultWebDeviceManager.InitWebDevicePool(); err != nil {
			log.Printf("Failed to initialize WebDeviceManager: %v", err)
		}
	})
	return nil
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
	m.mu.Lock()
	defer m.mu.Unlock()

	// 查找使用旧IP的设备
	if devices, ok := m.deviceIpMap[oldIp]; ok {
		// 更新这些设备的IP
		for _, dev := range devices {
			dev.ProxyIp = newIp
		}
		// 更新映射关系
		delete(m.deviceIpMap, oldIp)
		m.deviceIpMap[newIp] = devices
	}
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
	devices, err := m.webDeviceSvc.GetActiveByRangeId(currentIndex, currentIndex+pageSize)
	if err != nil {
		return fmt.Errorf("get web devices error: %v", err)
	}

	// 填充设备池
	for _, d := range devices {
		select {
		case m.webDevicePool <- d:
		default:
			// 设备池已满
			return nil
		}
	}

	return nil
}

// GetWebDevice 获取一个可用设备
func (m *WebDeviceManager) GetWebDevice() (*dto.WebDeviceDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), deviceTimeout)
	defer cancel()

	select {
	case dev := <-m.webDevicePool:
		// 检查设备是否达到使用上限
		key := fmt.Sprintf("%s_%d", m.webDeviceExpireKey, dev.Id)
		count, _ := strconv.ParseInt(redis.Get(key), 10, 64)
		maxUseConfig, err := m.dictionarySvc.GetByCode(constants.WEB_DEVICE_MAX_USE.Code)
		if err != nil {
			return nil, fmt.Errorf("get web device max use config error: %v", err)
		}
		maxUse, _ := strconv.ParseInt(maxUseConfig.Value, 10, 64)

		if count >= maxUse {
			m.mu.Lock()
			// 清空当前channel
			for len(m.webDevicePool) > 0 {
				<-m.webDevicePool
			}

			// 获取新的一批设备
			pageSize, err := m.getWebDeviceRange()
			if err != nil {
				m.mu.Unlock()
				return nil, err
			}

			currentIndex, err := m.getCurrentIndex()
			if err != nil {
				m.mu.Unlock()
				return nil, fmt.Errorf("get current index error: %v", err)
			}

			// 获取新的设备列表
			devices, err := m.webDeviceSvc.GetActiveByRangeId(currentIndex, currentIndex+pageSize)
			if err != nil {
				m.mu.Unlock()
				return nil, fmt.Errorf("get web devices error: %v", err)
			}

			// 如果没有新设备，重置索引并重新获取
			if len(devices) == 0 {
				minId, err := m.webDeviceSvc.MinIdByStartIndex(currentIndex)
				if err != nil {
					m.mu.Unlock()
					return nil, fmt.Errorf("get min id error: %v", err)
				}
				devices, err = m.webDeviceSvc.GetActiveByRangeId(minId, minId+pageSize)
				if err != nil {
					m.mu.Unlock()
					return nil, fmt.Errorf("get web devices error: %v", err)
				}
			}

			// 更新当前索引
			if len(devices) > 0 {
				err = m.updateCurrentIndex(int64(devices[len(devices)-1].Id))
				if err != nil {
					m.mu.Unlock()
					return nil, fmt.Errorf("update current index error: %v", err)
				}

				// 填充新设备到池中
				for _, d := range devices {
					m.webDevicePool <- d
				}
				m.mu.Unlock()

				// 返回第一个新设备
				return devices[0], nil
			}
			m.mu.Unlock()
			return nil, fmt.Errorf("no available web device")
		}

		// 更新设备使用次数
		_ = redis.Incr(key)
		redis.Expire(key, time.Duration(webDeviceExpireTime)*time.Second)

		// 设置默认值
		m.fillWebDeviceDefaultValues(dev)

		// 将设备放回池中
		go func() {
			select {
			case m.webDevicePool <- dev:
			case <-ctx.Done():
				log.Printf("放回设备池超时: %v", ctx.Err())
			}
		}()

		return dev, nil

	case <-ctx.Done():
		return nil, fmt.Errorf("获取设备超时: %v", ctx.Err())
	}
}

// fillWebDeviceDefaultValues 填充设备默认值
func (m *WebDeviceManager) fillWebDeviceDefaultValues(dev *dto.WebDeviceDTO) {
	if dev.Platform == "" {
		dev.Platform = "android"
	}
	if dev.UserAgent == "" {
		dev.UserAgent = ""
	}

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
	config, err := m.dictionarySvc.GetByCode("WEB_DEVICE_CURRENT_INDEX")
	if err != nil {
		return err
	}

	config.Value = strconv.FormatInt(index, 10)
	return m.dictionarySvc.Save(config)
}

// checkWebDeviceAvailable 检查设备是否可用
func (m *WebDeviceManager) checkWebDeviceAvailable(dev *dto.WebDeviceDTO) (bool, error) {
	key := fmt.Sprintf("%s_%d", m.webDeviceExpireKey, dev.Id)
	if !redis.Exists(key) {
		return true, nil
	}

	count, _ := strconv.ParseInt(redis.Get(key), 10, 64)

	// 获取设备最大使用次数
	maxUseConfig, err := m.dictionarySvc.GetByCode(constants.WEB_DEVICE_MAX_USE.Code)
	if err != nil {
		return false, fmt.Errorf("get web device max use config error: %v", err)
	}

	maxUse, _ := strconv.ParseInt(maxUseConfig.Value, 10, 64)
	return count < maxUse, nil
}

// getNewWebDevice 获取新设备
func (m *WebDeviceManager) getNewWebDevice() (*dto.WebDeviceDTO, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取设备范围
	pageSize, err := m.getWebDeviceRange()
	if err != nil {
		return nil, err
	}

	// 获取当前索引
	currentIndex, err := m.getCurrentIndex()
	if err != nil {
		return nil, fmt.Errorf("get current index error: %v", err)
	}

	// 获取新的设备列表
	devices, err := m.webDeviceSvc.GetActiveByRangeId(currentIndex, currentIndex+pageSize)
	if err != nil {
		return nil, fmt.Errorf("get web devices error: %v", err)
	}

	if len(devices) == 0 {
		// 如果没有新设备，重置索引
		minId, err := m.webDeviceSvc.MinIdByStartIndex(currentIndex)
		if err != nil {
			return nil, fmt.Errorf("get min id error: %v", err)
		}

		devices, err = m.webDeviceSvc.GetActiveByRangeId(minId, minId+pageSize)
		if err != nil {
			return nil, fmt.Errorf("get web devices error: %v", err)
		}
	}

	// 更新当前索引
	if len(devices) > 0 {
		err = m.updateCurrentIndex(int64(devices[len(devices)-1].Id))
		if err != nil {
			return nil, fmt.Errorf("update current index error: %v", err)
		}
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no available web device")
	}

	return devices[0], nil
}
