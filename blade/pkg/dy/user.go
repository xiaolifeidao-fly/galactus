package dy

import (
	webDeviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/dy"
	"galactus/blade/internal/service/dy/dto"
	"strconv"

	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

func UserRouters(engine *gin.RouterGroup) {
	engine.GET("/dy/user/get", getUserBySecUid)
	engine.GET("/dy/user/favorites", getUserFavoriteBySecUid)
	engine.GET("/dy/user/following", getUserFollowingBySecUid)

}

func getUserBySecUid(context *gin.Context) {
	secUid := context.Query("secUid")
	webDeviceDTO, _ := webDeviceService.FindWebDeviceById(6)
	userInfo := &dy.UserInfoEntity{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO),
		SecUid:       secUid,
	}
	result, err := dy.GetUserInfoByWeb(userInfo)
	routers.ToJson(context, result, err)
}

func getUserFavoriteBySecUid(context *gin.Context) {
	secUid := context.Query("secUid")
	maxCursor, _ := strconv.Atoi(context.Query("maxCursor"))
	minCursor, _ := strconv.Atoi(context.Query("minCursor"))
	webDeviceDTO, _ := webDeviceService.FindWebDeviceById(6)
	userFavoriteEntity := &dy.UserFavoriteEntity{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO),
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
	webDeviceDTO, _ := webDeviceService.FindWebDeviceById(6)
	userFollowingEntity := &dy.UserFollowingEntity{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO),
		SecUid:       secUid,
		Offset:       offset,
		UserId:       userId,
	}
	result, err := dy.GetUserFollowingByWeb(userFollowingEntity)
	routers.ToJson(context, result, err)
}
