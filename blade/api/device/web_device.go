package device

import (
	deviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/dto"
	"galactus/common/base/vo"
	"galactus/common/converter"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

type WebDevice struct {
	vo.Base
	DevicePlatform    string `json:"devicePlatform" description:"设备平台"`
	Aid               string `json:"aid" description:"应用ID"`
	Channel           string `json:"channel" description:"渠道"`
	Source            string `json:"source" description:"来源"`
	UpdateVersionCode string `json:"updateVersionCode" description:"更新版本代码"`
	PcClientType      string `json:"pcClientType" description:"PC客户端类型"`
	VersionCode       string `json:"versionCode" description:"版本代码"`
	VersionName       string `json:"versionName" description:"版本名称"`
	CookieEnabled     string `json:"cookieEnabled" description:"Cookie是否启用"`
	ScreenWidth       string `json:"screenWidth" description:"屏幕宽度"`
	ScreenHeight      string `json:"screenHeight" description:"屏幕高度"`
	BrowserLanguage   string `json:"browserLanguage" description:"浏览器语言"`
	BrowserPlatform   string `json:"browserPlatform" description:"浏览器平台"`
	BrowserName       string `json:"browserName" description:"浏览器名称"`
	BrowserVersion    string `json:"browserVersion" description:"浏览器版本"`
	BrowserOnline     string `json:"browserOnline" description:"浏览器是否在线"`
	EngineName        string `json:"engineName" description:"引擎名称"`
	EngineVersion     string `json:"engineVersion" description:"引擎版本"`
	OsName            string `json:"osName" description:"操作系统名称"`
	OsVersion         string `json:"osVersion" description:"操作系统版本"`
	CpuCoreNum        string `json:"cpuCoreNum" description:"CPU核心数"`
	DeviceMemory      string `json:"deviceMemory" description:"设备内存"`
	Platform          string `json:"platform" description:"平台"`
	Downlink          string `json:"downlink" description:"下载速度"`
	EffectiveType     string `json:"effectiveType" description:"有效类型"`
	RoundTripTime     string `json:"roundTripTime" description:"往返时间"`
	Webid             string `json:"webid" description:"WebID"`
}

func Routers(engine *gin.RouterGroup) {
	//用户列表
	engine.GET("/devices/list", list)
	engine.POST("/devices/save", save)

}

func save(context *gin.Context) {
	var device WebDevice
	context.ShouldBindJSON(&device)
	deviceDTO := converter.ToDTO[dto.WebDevice](&device)
	deviceDTO, err := deviceService.SaveWebDevice(deviceDTO)
	routers.ToJson(context, converter.ToVO[WebDevice](deviceDTO), err)
}

func list(context *gin.Context) {
	devices, err := deviceService.ListWebDevice()
	routers.ToJson(context, converter.ToVOs[WebDevice](devices), err)
}
