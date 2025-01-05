package dy

import (
	webDeviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/dy"
	"galactus/blade/internal/service/dy/dto"
	"strconv"

	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

// TODO: 下面获取webDeviceDTO的逻辑要抽取出来一个模块，从里面获取一个设备，每次尽可能获取不一样的设备，轮询使用
// TODO: web_device 表增加IP字段, 一个web_device 对应一个IP，这个IP失效的话，重新更新
// TODO: 获取IP的逻辑要抽取出来一个模块，从里面获取一个IP，每次尽可能获取不一样的IP，轮询使用

func UserRouters(engine *gin.RouterGroup) {
	engine.GET("/dy/user/get", getUserBySecUid)
	engine.GET("/dy/user/favorites", getUserFavoriteBySecUid)
	engine.GET("/dy/user/following", getUserFollowingBySecUid)
	engine.GET("/dy/user/convertUidByUrl", convertUidByUrl)

}

func convertUidByUrl(context *gin.Context) {
	url := context.Query("url")
	webDeviceDTO, _ := webDeviceService.NewWebDeviceService().GetById(6)
	ip := "" //TODO 获取IP
	userInfoEntity := &dy.UserInfoEntity{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
	}
	result := dy.GetUrlByUrl(url, userInfoEntity)
	routers.ToJson(context, result, nil)
}

func getUserBySecUid(context *gin.Context) {
	secUid := context.Query("secUid")
	webDeviceDTO, _ := webDeviceService.NewWebDeviceService().GetById(6)
	ip := "" //TODO 获取IP
	userInfo := &dy.UserInfoEntity{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
		SecUid:       secUid,
	}
	result, err := dy.GetUserInfoByWeb(userInfo)
	routers.ToJson(context, result, err)
}

func getUserFavoriteBySecUid(context *gin.Context) {
	secUid := context.Query("secUid")
	maxCursor, _ := strconv.Atoi(context.Query("maxCursor"))
	minCursor, _ := strconv.Atoi(context.Query("minCursor"))
	webDeviceDTO, _ := webDeviceService.NewWebDeviceService().GetById(6)
	ip := "" //TODO 获取IP
	userFavoriteEntity := &dy.UserFavoriteEntity{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
		SecUid:       secUid,
		MaxCursor:    maxCursor,
		MinCursor:    minCursor,
	}
	result, err := dy.GetUserFavoriteByWeb(userFavoriteEntity)
	routers.ToJson(context, result, err)
}

func getUserFollowingBySecUid(context *gin.Context) {
	secUid := context.Query("secUid")
	offset, _ := strconv.Atoi(context.Query("offset"))
	userId := context.Query("userId")
	webDeviceDTO, _ := webDeviceService.NewWebDeviceService().GetById(6)
	ip := "" //TODO 获取IP
	userFollowingEntity := &dy.UserFollowingEntity{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
		SecUid:       secUid,
		Offset:       offset,
		UserId:       userId,
	}
	result, err := dy.GetUserFollowingByWeb(userFollowingEntity)
	routers.ToJson(context, result, err)
}
