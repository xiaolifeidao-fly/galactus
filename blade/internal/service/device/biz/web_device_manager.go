package biz

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/dictionary"
	"galactus/blade/internal/service/ip/biz"
	"galactus/common/middleware/redis"
)

const (
	webDeviceExpireTime = 30 * 60 // 30分钟
	webDevicePageSize   = 5000
	webDevicePoolKey    = "WEB_DEVICE_POOL"
)

type WebDeviceManager struct {
	mu                 sync.Mutex
	webDevicePool      chan *dto.WebDeviceDTO
	webDeviceSvc       *device.WebDeviceService
	dictionarySvc      dictionary.DictionaryService
	webDeviceExpireKey string
	ipManager          *biz.IpManager
	deviceIpMap        map[string][]*dto.WebDeviceDTO // key: ip, value: devices using this ip
}

func NewWebDeviceManager(
	webDeviceSvc *device.WebDeviceService,
	dictionarySvc dictionary.DictionaryService,
) *WebDeviceManager {
	manager := &WebDeviceManager{
		webDevicePool:      make(chan *dto.WebDeviceDTO, webDevicePageSize),
		webDeviceSvc:       webDeviceSvc,
		dictionarySvc:      dictionarySvc,
		webDeviceExpireKey: "WEB_DEVICE_EXPIRE",
		ipManager:          biz.GetDefaultIpManager(),
		deviceIpMap:        make(map[string][]*dto.WebDeviceDTO),
	}
	// 注册为IP更新的观察者
	manager.ipManager.RegisterObserver(manager)
	return manager
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
func (m *WebDeviceManager) InitWebDevicePool(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取设备ID范围
	_, err := m.dictionarySvc.GetByCode("WEB_DEVICE_ID_RANGE")
	if err != nil {
		return fmt.Errorf("get web device range config error: %v", err)
	}

	// 获取当前索引
	currentIndex, err := m.getCurrentIndex(ctx)
	if err != nil {
		return fmt.Errorf("get current index error: %v", err)
	}

	// 获取设备列表
	devices, err := m.webDeviceSvc.GetActiveByRangeId(currentIndex, currentIndex+webDevicePageSize)
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
func (m *WebDeviceManager) GetWebDevice(ctx context.Context) (*dto.WebDeviceDTO, error) {
	select {
	case dev := <-m.webDevicePool:
		// 检查设备是否可用
		if ok, err := m.checkWebDeviceAvailable(ctx, dev); err != nil {
			return nil, err
		} else if !ok {
			// 设备不可用，尝试获取新设备
			return m.getNewWebDevice(ctx)
		}

		// 更新设备使用次数
		key := fmt.Sprintf("%s_%d", m.webDeviceExpireKey, dev.Id)
		_ = redis.Incr(key)

		// 设置过期时间
		redis.Expire(key, time.Duration(webDeviceExpireTime)*time.Second)

		// 设置默认值
		m.fillWebDeviceDefaultValues(dev)

		// 将设备放回池中
		go func() {
			m.webDevicePool <- dev
		}()

		return dev, nil

	case <-ctx.Done():
		return nil, ctx.Err()
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

// checkWebDeviceAvailable 检查设备是否可用
func (m *WebDeviceManager) checkWebDeviceAvailable(ctx context.Context, dev *dto.WebDeviceDTO) (bool, error) {
	key := fmt.Sprintf("%s_%d", m.webDeviceExpireKey, dev.Id)
	if !redis.Exists(key) {
		return true, nil
	}

	count, _ := strconv.ParseInt(redis.Get(key), 10, 64)

	// 获取设备最大使用次数
	maxUseConfig, err := m.dictionarySvc.GetByCode("WEB_DEVICE_MAX_USE")
	if err != nil {
		return false, fmt.Errorf("get web device max use config error: %v", err)
	}

	maxUse, _ := strconv.ParseInt(maxUseConfig.Value, 10, 64)
	return count < maxUse, nil
}

// getNewWebDevice 获取新设备
func (m *WebDeviceManager) getNewWebDevice(ctx context.Context) (*dto.WebDeviceDTO, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取当前索引
	currentIndex, err := m.getCurrentIndex(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current index error: %v", err)
	}

	// 获取新的设备列表
	devices, err := m.webDeviceSvc.GetActiveByRangeId(currentIndex, currentIndex+webDevicePageSize)
	if err != nil {
		return nil, fmt.Errorf("get web devices error: %v", err)
	}

	if len(devices) == 0 {
		// 如果没有新设备，重置索引
		minId, err := m.webDeviceSvc.MinIdByStartIndex(currentIndex)
		if err != nil {
			return nil, fmt.Errorf("get min id error: %v", err)
		}

		devices, err = m.webDeviceSvc.GetActiveByRangeId(minId, minId+webDevicePageSize)
		if err != nil {
			return nil, fmt.Errorf("get web devices error: %v", err)
		}
	}

	// 更新当前索引
	if len(devices) > 0 {
		err = m.updateCurrentIndex(ctx, int64(devices[len(devices)-1].Id))
		if err != nil {
			return nil, fmt.Errorf("update current index error: %v", err)
		}
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no available web device")
	}

	return devices[0], nil
}

// getCurrentIndex 获取当前索引
func (m *WebDeviceManager) getCurrentIndex(ctx context.Context) (int64, error) {
	config, err := m.dictionarySvc.GetByCode("WEB_DEVICE_CURRENT_INDEX")
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(config.Value, 10, 64)
}

// updateCurrentIndex 更新当前索引
func (m *WebDeviceManager) updateCurrentIndex(ctx context.Context, index int64) error {
	config, err := m.dictionarySvc.GetByCode("WEB_DEVICE_CURRENT_INDEX")
	if err != nil {
		return err
	}

	config.Value = strconv.FormatInt(index, 10)
	return m.dictionarySvc.Save(config)
}
