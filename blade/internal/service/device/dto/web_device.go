package dto

import (
	"galactus/common/base/dto"
	"time"
)

type WebDeviceDTO struct {
	dto.BaseDTO
	DevicePlatform    string    `json:"device_platform"`
	Aid               string    `json:"aid"`
	Channel           string    `json:"channel"`
	Source            string    `json:"source"`
	UpdateVersionCode string    `json:"update_version_code"`
	PcClientType      string    `json:"pc_client_type"`
	VersionCode       string    `json:"version_code"`
	VersionName       string    `json:"version_name"`
	CookieEnabled     string    `json:"cookie_enabled"`
	ScreenWidth       string    `json:"screen_width"`
	ScreenHeight      string    `json:"screen_height"`
	BrowserLanguage   string    `json:"browser_language"`
	BrowserPlatform   string    `json:"browser_platform"`
	BrowserName       string    `json:"browser_name"`
	BrowserVersion    string    `json:"browser_version"`
	BrowserOnline     string    `json:"browser_online"`
	EngineName        string    `json:"engine_name"`
	EngineVersion     string    `json:"engine_version"`
	OsName            string    `json:"os_name"`
	OsVersion         string    `json:"os_version"`
	CpuCoreNum        string    `json:"cpu_core_num"`
	DeviceMemory      string    `json:"device_memory"`
	Platform          string    `json:"platform"`
	Downlink          string    `json:"downlink"`
	EffectiveType     string    `json:"effective_type"`
	RoundTripTime     string    `json:"round_trip_time"`
	Webid             string    `json:"webid"`
	Uifid             string    `json:"uifid"`
	VerifyFp          string    `json:"verify_fp"`
	Fp                string    `json:"fp"`
	Ttwid             string    `json:"ttwid"`
	OdinTt            string    `json:"odin_tt"`
	UserAgent         string    `json:"user_agent"`
	ProxyIp           string    `json:"proxy_ip"`
	Cookie            string    `json:"cookie"`
	PcLibraDivert     string    `json:"pc_libra_divert"`
	SecChUaPlatform   string    `json:"sec_ch_ua_platform"`
	SecChUa           string    `json:"sec_ch_ua"`
	ExpireTime        time.Time `json:"expire_time"`
}
