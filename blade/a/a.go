package a

import (
	"galactus/blade/internal/service/ip"
	"galactus/blade/internal/service/ip/biz"
	"galactus/common/middleware/db"
	"galactus/common/middleware/redis"
	"galactus/common/middleware/vipper"
	"log"
)

func init() {
	vipper.Init()
}

func Init() {
	// db init
	db.InitDB()

	// redis init
	redisAddr := vipper.GetString("redis.addr")
	redisPwd := vipper.GetString("redis.password")
	redis.InitRedisClient(redisAddr, redisPwd)

	// zdy http proxy service init
	zdyUrl := vipper.GetString("zdy.url")
	zdyKey := vipper.GetString("zdy.key")
	zdyApi := vipper.GetString("zdy.api")
	ip.InitDefaultZDYHttpProxyService(zdyUrl, zdyKey, zdyApi)

	// ip service init
	biz.InitDefaultIpManager()
	err := biz.GetDefaultIpManager().InitIp()
	if err != nil {
		log.Printf("ip init failed: %v", err)
	}
}
