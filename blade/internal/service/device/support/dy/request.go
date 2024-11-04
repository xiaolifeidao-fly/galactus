package dy

import (
	"fmt"
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/device/support"
	"galactus/common/middleware/http"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

var singUrl = "http://localhost:9001"

type DyEntity interface {
	GetAbogus(params string, ua string) string
	GetAcSign(url string, acNonce string, ua string) string
}

type DyBaseEntity struct {
	support.Entity
	webDevice *dto.WebDeviceDTO
	url       string
	method    string
	body      map[string]interface{}
}

func (e *DyBaseEntity) GetCookie() map[string]interface{} {
	return map[string]interface{}{
		"odin_tt": e.webDevice.OdinTt,
		"ttwid":   e.webDevice.Ttwid,
		"UIFID":   e.webDevice.Uifid,
	}
}

func newDyBaseEntity(url string, method string, webDevice *dto.WebDeviceDTO) *DyBaseEntity {
	dyBaseEntity := &DyBaseEntity{
		url:       url,
		method:    method,
		webDevice: webDevice,
	}
	dyBaseEntity.AppendCommonParams()
	return dyBaseEntity
}

func (e *DyBaseEntity) GetCookieString() string {
	cookie := e.GetCookie()
	cookieString := ""
	for key, value := range cookie {
		cookieString += fmt.Sprintf("%s=%s;", key, value)
	}
	return cookieString
	// return "ttwid=1%7ChvqFO57kqn8BATFRUUW8eXgPhICpkh92LVPb4lGCyhU%7C1729751805%7Cc3a206e2a10b9760da63720198d3ca1be1187a1e90280f6294aa58e4e288fe81; "
}

func (e *DyBaseEntity) GetMethod() string {
	return e.method
}

func (e *DyBaseEntity) GetBody() map[string]interface{} {
	return e.body
}

func (e *DyBaseEntity) GetHeaders() map[string]string {
	headers := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "zh-CN,zh;q=0.9",
		"origin":             "https://www.douyin.com",
		"priority":           "u=1, i",
		"referer":            "https://www.douyin.com/",
		"sec-ch-ua":          "Chromium;v=130, Google Chrome;v=130, Not?A_Brand;v=99",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "macOS",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"uifid":              e.webDevice.Uifid,
		"user-agent":         e.webDevice.UserAgent,
	}
	return headers
}

func (e *DyBaseEntity) GetCommonParams() map[string]interface{} {
	return map[string]interface{}{
		"device_platform":     e.webDevice.DevicePlatform,
		"aid":                 e.webDevice.Aid,
		"channel":             e.webDevice.Channel,
		"source":              e.webDevice.Source,
		"update_version_code": e.webDevice.UpdateVersionCode,
		"pc_client_type":      e.webDevice.PcClientType,
		"version_code":        e.webDevice.VersionCode,
		"version_name":        e.webDevice.VersionName,
		"cookie_enabled":      e.webDevice.CookieEnabled,
		"screen_width":        e.webDevice.ScreenWidth,
		"screen_height":       e.webDevice.ScreenHeight,
		"browser_language":    e.webDevice.BrowserLanguage,
		"browser_platform":    e.webDevice.BrowserPlatform,
		"browser_name":        e.webDevice.BrowserName,
		"browser_version":     e.webDevice.BrowserVersion,
		"browser_online":      e.webDevice.BrowserOnline,
		"engine_name":         e.webDevice.EngineName,
		"engine_version":      e.webDevice.EngineVersion,
		"os_name":             url.QueryEscape(e.webDevice.OsName),
		"os_version":          e.webDevice.OsVersion,
		"cpu_core_num":        e.webDevice.CpuCoreNum,
		"device_memory":       e.webDevice.DeviceMemory,
		"platform":            e.webDevice.Platform,
		"downlink":            e.webDevice.Downlink,
		"effective_type":      e.webDevice.EffectiveType,
		"round_trip_time":     e.webDevice.RoundTripTime,
		"webid":               e.webDevice.Webid,
		"uifid":               e.webDevice.Uifid,
		"verify_fp":           e.webDevice.VerifyFp,
		"fp":                  e.webDevice.Fp,
		"msToken":             e.getMsToken(107),
	}
}

func (e *DyBaseEntity) GetParams() string {
	split := strings.Split(e.url, "?")
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

func (e *DyBaseEntity) AppendCommonParams() {
	params := e.GetCommonParams()
	for key, value := range params {
		e.AppendUrlParams(key, value.(string))
	}
}

func (e *DyBaseEntity) AppendUrlParams(name string, value string) {
	if e.url[len(e.url)-1] != '?' {
		e.url += "&" + name + "=" + value
	} else {
		e.url += name + "=" + value
	}
}

func (e *DyBaseEntity) getMsToken(randomLength int) string {
	// 根据传入长度产生随机字符串
	baseStr := "ABCDEFGHIGKLMNOPQRSTUVWXYZabcdefghigklmnopqrstuvwxyz0123456789="
	length := len(baseStr)
	randomStr := make([]byte, randomLength)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < randomLength; i++ {
		randomStr[i] = baseStr[rand.Intn(length)]
	}

	return string(randomStr)
}

func (e *DyBaseEntity) GetAbogus(params string, ua string) string {
	result, _ := http.Post(singUrl+"/dy/abogus/sign", map[string]interface{}{
		"params": params,
		"ua":     ua,
	}, "", nil)
	return result["aBogus"].(string)
}

func (e *DyBaseEntity) GetAcSign(url string, acNonce string, ua string) string {
	result, _ := http.Post(singUrl+"/dy/ac/sign", map[string]interface{}{
		"url":     url,
		"acNonce": acNonce,
		"ua":      ua,
	}, "", nil)
	return result["acSignature"].(string)
}

func (e *DyBaseEntity) GetUrl() string {
	return e.url
}

type DyRequest[E support.Entity] struct {
	support.Request[E]
}

func (r *DyRequest[E]) DoRequest(e E) (map[string]interface{}, error) {
	if e.GetMethod() == "GET" {
		result, err := http.Get(e.GetUrl(), e.GetCookieString(), e.GetHeaders())
		return result, err
	}
	result, err := http.Post(e.GetUrl(), e.GetBody(), e.GetCookieString(), e.GetHeaders())
	return result, err
}

func DoRequest(e support.Entity) (map[string]interface{}, error) {
	requestInstance := &DyRequest[support.Entity]{}
	return requestInstance.DoRequest(e)
}
