package a

import (
	"galactus/blade/internal/service/ip"
	"galactus/common/middleware/db"
	"galactus/common/middleware/redis"
	"galactus/common/middleware/vipper"
	"log"
)

func init() {
	vipper.Init()
}

func Init() {

	// 再初始化数据库
	db.InitDB()

	// redis init
	redisAddr := vipper.GetString("redis.addr")
	redisPwd := vipper.GetString("redis.password")
	redis.InitRedisClient(redisAddr, redisPwd)

	// ip service init
	_, err := ip.InitDefaultIpService()
	if err != nil {
		log.Printf("ip init failed: %v", err)
	}

	// zdy http proxy service init
	zdyUrl := vipper.GetString("zdy.url")
	zdyKey := vipper.GetString("zdy.key")
	zdyApi := vipper.GetString("zdy.api")
	ip.InitDefaultZDYHttpProxyService(zdyUrl, zdyKey, zdyApi)
}
