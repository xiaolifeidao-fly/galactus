package initialize

import (
	"fmt"
	deviceBiz "galactus/blade/internal/service/device/biz"
	"galactus/blade/internal/service/ip"
	ipBiz "galactus/blade/internal/service/ip/biz"
	"galactus/common/middleware/db"
	"galactus/common/middleware/redis"
	"galactus/common/middleware/vipper"
	"log"
)

// InitOrder 定义初始化顺序
type InitOrder int

const (
	ConfigInit InitOrder = iota
	DBInit
	RedisInit
	ProxyServiceInit
	IpManagerInit
	DeviceManagerInit
)

type Initializer struct {
	Order  InitOrder
	Name   string
	InitFn func() error
}

var initializers = []Initializer{
	{
		Order: ConfigInit,
		Name:  "Config",
		InitFn: func() error {
			vipper.Init()
			return nil
		},
	},
	{
		Order: DBInit,
		Name:  "Database",
		InitFn: func() error {
			db.InitDB()
			return nil
		},
	},
	{
		Order: RedisInit,
		Name:  "Redis",
		InitFn: func() error {
			redisAddr := vipper.GetString("redis.addr")
			redisPwd := vipper.GetString("redis.password")
			return redis.InitRedisClient(redisAddr, redisPwd)
		},
	},
	{
		Order: ProxyServiceInit,
		Name:  "ZDY Proxy Service",
		InitFn: func() error {
			zdyUrl := vipper.GetString("zdy.url")
			zdyKey := vipper.GetString("zdy.key")
			zdyApi := vipper.GetString("zdy.api")
			return ip.InitDefaultZDYHttpProxyService(zdyUrl, zdyKey, zdyApi)
		},
	},
	{
		Order:  IpManagerInit,
		Name:   "IP Manager",
		InitFn: ipBiz.InitIpManager,
	},
	{
		Order:  DeviceManagerInit,
		Name:   "Device Manager",
		InitFn: deviceBiz.InitDefaultWebDeviceManager,
	},
}

// Initialize 统一初始化入口
func Initialize() error {
	// 按顺序执行初始化
	for _, init := range initializers {
		log.Printf("Initializing %s...", init.Name)
		if err := init.InitFn(); err != nil {
			return fmt.Errorf("failed to initialize %s: %v", init.Name, err)
		}
		log.Printf("%s initialized successfully", init.Name)
	}
	return nil
}
