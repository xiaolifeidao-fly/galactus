package dy

import (
	"errors"
	"galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/biz"
	"galactus/blade/internal/service/dy"
	"galactus/blade/internal/service/dy/dto"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

// TODO: 下面获取webDeviceDTO的逻辑要抽取出来一个模块，从里面获取一个设备，每次尽可能获取不一样的设备，轮询使用
// TODO: web_device 表增加IP字段, 一个web_device 对应一个IP，这个IP失效的话，重新更新
// TODO: 获取IP的逻辑要抽取出来一个模块，从里面获取一个IP，每次尽可能获取不一样的IP，轮询使用

type VideoHandler struct {
	*routers.BaseHandler
	*biz.WebDeviceManager
	WebDeviceService *device.WebDeviceService
}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{
		WebDeviceManager: biz.GetDefaultWebDeviceManager(),
		WebDeviceService: device.NewWebDeviceService(),
	}
}

func (h *VideoHandler) RegisterHandler(engine *gin.RouterGroup) {
	engine.GET("/dy/video/player", h.playerVideo)
	engine.GET("/dy/video/love", h.loveVideo)
	engine.GET("/dy/video/convert", h.convertByVideoUrl)
	engine.GET("/dy/video/get", h.getVideoInfo)
}

func (h *VideoHandler) convertByVideoUrl(context *gin.Context) {
	businessKey := context.Query("businessKey")
	ip := "" //TODO 获取IP
	response := dy.ConvertByVideoUrl(businessKey, ip)
	routers.ToJson(context, response, nil)
}

func (h *VideoHandler) getVideoInfo(context *gin.Context) {
	videoId := context.Query("videoId")
	webDeviceDTO, _ := h.GetWebDevice()
	if webDeviceDTO == nil {
		routers.ToJson(context, nil, errors.New("webDeviceDTO is nil"))
		return
	}
	ip := "" //TODO 获取IP
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
		VideoId:      videoId,
	}
	result := dy.GetVideoItemInfo(videoInfo)
	routers.ToJson(context, result, nil)
}

func (h *VideoHandler) playerVideo(context *gin.Context) {
	videoId := context.Query("videoId")
	webDeviceDTO, _ := h.GetWebDevice()
	ip := "" //TODO 获取IP
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
		VideoId:      videoId,
	}
	result, err := dy.PlayerVideo(videoInfo)
	routers.ToJson(context, result, err)
}

func (h *VideoHandler) loveVideo(context *gin.Context) {
	videoId := context.Query("videoId")
	webDeviceDTO, _ := h.GetWebDevice()
	ip := "" //TODO 获取IP
	videoInfo := &dy.VideoInfo{
		DyBaseEntity: dto.NewDyBaseEntity(webDeviceDTO, ip),
		VideoId:      videoId,
	}
	result, err := dy.LoveVideo(videoInfo)
	routers.ToJson(context, result, err)
}
