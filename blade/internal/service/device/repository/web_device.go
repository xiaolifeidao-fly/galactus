package repository

import (
	"galactus/common/middleware/db"
)

type WebDevice struct {
	db.BaseEntity
	DevicePlatform    string `orm:"column(device_platform);size(32);null" description:"设备平台"`
	Aid               string `orm:"column(aid);size(32);null" description:"应用ID"`
	Channel           string `orm:"column(channel);size(32);null" description:"渠道"`
	Source            string `orm:"column(source);size(32);null" description:"来源"`
	UpdateVersionCode string `orm:"column(update_version_code);size(32);null" description:"更新版本代码"`
	PcClientType      string `orm:"column(pc_client_type);size(32);null" description:"PC客户端类型"`
	VersionCode       string `orm:"column(version_code);size(32);null" description:"版本代码"`
	VersionName       string `orm:"column(version_name);size(32);null" description:"版本名称"`
	CookieEnabled     string `orm:"column(cookie_enabled);size(32);null" description:"Cookie是否启用"`
	ScreenWidth       string `orm:"column(screen_width);size(32);null" description:"屏幕宽度"`
	ScreenHeight      string `orm:"column(screen_height);size(32);null" description:"屏幕高度"`
	BrowserLanguage   string `orm:"column(browser_language);size(32);null" description:"浏览器语言"`
	BrowserPlatform   string `orm:"column(browser_platform);size(32);null" description:"浏览器平台"`
	BrowserName       string `orm:"column(browser_name);size(32);null" description:"浏览器名称"`
	BrowserVersion    string `orm:"column(browser_version);size(32);null" description:"浏览器版本"`
	BrowserOnline     string `orm:"column(browser_online);size(32);null" description:"浏览器是否在线"`
	EngineName        string `orm:"column(engine_name);size(32);null" description:"引擎名称"`
	EngineVersion     string `orm:"column(engine_version);size(32);null" description:"引擎版本"`
	OsName            string `orm:"column(os_name);size(32);null" description:"操作系统名称"`
	OsVersion         string `orm:"column(os_version);size(32);null" description:"操作系统版本"`
	CpuCoreNum        string `orm:"column(cpu_core_num);size(32);null" description:"CPU核心数"`
	DeviceMemory      string `orm:"column(device_memory);size(32);null" description:"设备内存"`
	Platform          string `orm:"column(platform);size(32);null" description:"平台"`
	Downlink          string `orm:"column(downlink);size(32);null" description:"下载速度"`
	EffectiveType     string `orm:"column(effective_type);size(32);null" description:"有效类型"`
	RoundTripTime     string `orm:"column(round_trip_time);size(32);null" description:"往返时间"`
	Webid             string `orm:"column(webid);size(32);null" description:"WebID"`
	Uifid             string `orm:"column(uifid);size(1000);null" description:"UIFID"`
	VerifyFp          string `orm:"column(verify_fp);size(500);null" description:"验证FP"`
	Fp                string `orm:"column(fp);size(500);null" description:"FP"`
	Ttwid             string `orm:"column(ttwid);size(500);null" description:"TTWID"`
	OdinTt            string `orm:"column(odin_tt);size(500);null" description:"ODINTT"`
	UserAgent         string `orm:"column(user_agent);size(500);null" description:"UserAgent"`
}

func (w *WebDevice) TableName() string {
	return "web_device"
}

type WebDeviceRepository struct {
	db.Repository[*WebDevice]
}

func (w *WebDeviceRepository) GetByWebid(webid string) (*WebDevice, error) {
	device, err := w.GetOne("select * from web_device where webid = ?", webid)
	return device, err
}
