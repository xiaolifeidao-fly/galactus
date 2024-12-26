package dy

import (
	webDeviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/dy"
	"galactus/blade/internal/service/dy/dto"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

func VideoRouters(engine *gin.RouterGroup) {
	engine.GET("/dy/video/get", getVideoInfo)
	engine.GET("/dy/video/player", playerVideo)
	engine.GET("/dy/video/love", loveVideo)
}

func getVideoInfo(context *gin.Context) {
	videoId := context.Query("videoId")
	webDeviceDTO, _ := webDeviceService.FindWebDeviceById(6)
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO),
		VideoId:      videoId,
	}
	result, err := dy.GetVideoInfo(videoInfo)
	routers.ToJson(context, result, err)
}

func playerVideo(context *gin.Context) {
	videoId := context.Query("videoId")
	webDeviceDTO, _ := webDeviceService.FindWebDeviceById(6)
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO),
		VideoId:      videoId,
	}
	result, err := dy.PlayerVideo(videoInfo)
	routers.ToJson(context, result, err)
}

func loveVideo(context *gin.Context) {
	videoId := context.Query("videoId")
	webDeviceDTO, _ := webDeviceService.FindWebDeviceById(6)
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO),
		VideoId:      videoId,
	}
	result, err := dy.LoveVideo(videoInfo)
	routers.ToJson(context, result, err)
}
