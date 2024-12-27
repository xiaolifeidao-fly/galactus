package dy

import (
	webDeviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/dy"
	"galactus/blade/internal/service/dy/dto"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

func VideoRouters(engine *gin.RouterGroup) {
	engine.GET("/dy/video/player", playerVideo)
	engine.GET("/dy/video/love", loveVideo)
	engine.GET("/dy/video/convert", convertByVideoUrl)
	engine.GET("/dy/video/get", getVideoInfo)

}

func convertByVideoUrl(context *gin.Context) {
	businessKey := context.Query("businessKey")
	response := dy.ConvertByVideoUrl(businessKey)
	routers.ToJson(context, response, nil)
}

func getVideoInfo(context *gin.Context) {
	videoId := context.Query("videoId")
	webDeviceDTO, _ := webDeviceService.FindWebDeviceById(6)
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO),
		VideoId:      videoId,
	}
	result := dy.GetVideoItemInfo(videoInfo)
	routers.ToJson(context, result, nil)
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
