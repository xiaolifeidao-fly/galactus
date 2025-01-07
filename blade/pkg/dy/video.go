package dy

import (
	webDeviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/biz"
	"galactus/blade/internal/service/dy"
	"galactus/blade/internal/service/dy/dto"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

// TODO: 下面获取webDeviceDTO的逻辑要抽取出来一个模块，从里面获取一个设备，每次尽可能获取不一样的设备，轮询使用
// TODO: web_device 表增加IP字段, 一个web_device 对应一个IP，这个IP失效的话，重新更新
// TODO: 获取IP的逻辑要抽取出来一个模块，从里面获取一个IP，每次尽可能获取不一样的IP，轮询使用

var deviceManager = biz.GetDefaultWebDeviceManager()

func VideoRouters(engine *gin.RouterGroup) {
	engine.GET("/dy/video/player", playerVideo)
	engine.GET("/dy/video/love", loveVideo)
	engine.GET("/dy/video/convert", convertByVideoUrl)
	engine.GET("/dy/video/get", getVideoInfo)

}

func convertByVideoUrl(context *gin.Context) {
	businessKey := context.Query("businessKey")
	ip := "" //TODO 获取IP
	response := dy.ConvertByVideoUrl(businessKey, ip)
	routers.ToJson(context, response, nil)
}

func getVideoInfo(context *gin.Context) {
	videoId := context.Query("videoId");
	deviceManager.GetWebDevice(context.Request.Context())

	webDeviceDTO, _ := webDeviceService.NewWebDeviceService().GetById(15)
	ip := "" //TODO 获取IP
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
		VideoId:      videoId,
	}
	result := dy.GetVideoItemInfo(videoInfo)
	routers.ToJson(context, result, nil)
}

func playerVideo(context *gin.Context) {
	videoId := context.Query("videoId")
	webDeviceDTO, _ := webDeviceService.NewWebDeviceService().GetById(6)
	ip := "" //TODO 获取IP
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
		VideoId:      videoId,
	}
	result, err := dy.PlayerVideo(videoInfo)
	routers.ToJson(context, result, err)
}

func loveVideo(context *gin.Context) {
	videoId := context.Query("videoId")
	webDeviceDTO, _ := webDeviceService.NewWebDeviceService().GetById(31)
	ip := "" //TODO 获取IP
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
		VideoId:      videoId,
	}
	result, err := dy.LoveVideo(videoInfo)
	routers.ToJson(context, result, err)
}
