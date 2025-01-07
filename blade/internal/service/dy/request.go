package dy

import (
	"fmt"
	"galactus/blade/internal/service"
	"galactus/blade/internal/service/device/dto"
	"galactus/common/middleware/http"
	"galactus/common/utils"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type DyEntity interface {
	Init(Url string)
	GetAbogus(params string, ua string) string
	GetAcSign(Url string, acNonce string, ua string) string
}

type DyBaseEntity struct {
	service.Entity
	WebDevice *dto.WebDeviceDTO
	Ip        string
	Url       string
	Method    string
	Body      map[string]interface{}
}

func (e *DyBaseEntity) GetCookie() map[string]interface{} {
	return map[string]interface{}{
		"odin_tt": e.WebDevice.OdinTt,
		"ttwid":   e.WebDevice.Ttwid,
		"UIFID":   e.WebDevice.Uifid,
	}
}

func (e *DyBaseEntity) GetIp() string {
	return e.Ip
}

func (e *DyBaseEntity) Init(url string) {
	e.Url = url
	e.AppendCommonParams()
}

func (e *DyBaseEntity) GetCookieString() string {
	cookie := e.GetCookie()
	cookieString := ""
	for key, value := range cookie {
		cookieString += fmt.Sprintf("%s=%s;", key, value)
	}
	return cookieString
}

func (e *DyBaseEntity) GetMethod() string {
	return e.Method
}

func (e *DyBaseEntity) GetBody() map[string]interface{} {
	return e.Body
}

func (e *DyBaseEntity) GetHeaders() map[string]string {
	/**
	  -H 'bd-ticket-guard-client-data: eyJ0c19zaWduIjoidHMuMS40NTg4MzQ3MTcyYzhhYmJjZWZmZWFhMTBhNTg2YjQwM2QyZTY2OWY3YTQ2MWYwMjc0YjJiZTlmN2Y4MTQwNTNhYzRmYmU4N2QyMzE5Y2YwNTMxODYyNGNlZGExNDkxMWNhNDA2ZGVkYmViZWRkYjJlMzBmY2U4ZDRmYTAyNTc1ZCIsInJlcV9jb250ZW50IjoidGlja2V0LHBhdGgsdGltZXN0YW1wIiwicmVxX3NpZ24iOiJNRVVDSUVwSXVhMmJENWZKc0N0RHZsOVdKZEluQUJGTkdEUTlQS1kxdFZDRld6bGJBaUVBemo2cjF0ZHhBMWZYeThqUitnODlmREhQTC83dEpFbnkwZm1TYmhXVjN1VT0iLCJ0aW1lc3RhbXAiOjE3MzA3NzQ3ODd9' \
	  -H 'bd-ticket-guard-iteration-version: 1' \
	  -H 'bd-ticket-guard-ree-public-key: BDer+v4VrT/2TXp4LxgMhGwh20ikdwblB7luFglJabpT3fz8lshbB4AUiNTNuu1VC1A3Y7p6xQa//5hszKL3LVg=' \
	  -H 'bd-ticket-guard-version: 2' \
	  -H 'bd-ticket-guard-web-version: 1' \
	*/
	headers := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "zh-CN,zh;q=0.9",
		"origin":             "https://www.douyin.com",
		"priority":           "u=1, i",
		"referer":            "https://www.douyin.com/",
		"sec-ch-ua":          e.WebDevice.SecChUa,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": e.WebDevice.SecChUaPlatform,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"host":               "www.douyin.com",
		"uifid":              e.WebDevice.Uifid,
		"user-agent":         e.WebDevice.UserAgent,
	}
	return headers
}

func (e *DyBaseEntity) GetCommonParams() map[string]interface{} {
	return map[string]interface{}{
		"device_platform":     e.WebDevice.DevicePlatform,
		"aid":                 e.WebDevice.Aid,
		"channel":             e.WebDevice.Channel,
		"source":              e.WebDevice.Source,
		"update_version_code": e.WebDevice.UpdateVersionCode,
		"pc_client_type":      e.WebDevice.PcClientType,
		"pc_libra_divert":     e.WebDevice.PcLibraDivert,
		"version_code":        e.WebDevice.VersionCode,
		"version_name":        e.WebDevice.VersionName,
		"cookie_enabled":      e.WebDevice.CookieEnabled,
		"screen_width":        e.WebDevice.ScreenWidth,
		"screen_height":       e.WebDevice.ScreenHeight,
		"browser_language":    e.WebDevice.BrowserLanguage,
		"browser_platform":    e.WebDevice.BrowserPlatform,
		"browser_name":        e.WebDevice.BrowserName,
		"browser_version":     e.WebDevice.BrowserVersion,
		"browser_online":      e.WebDevice.BrowserOnline,
		"engine_name":         e.WebDevice.EngineName,
		"engine_version":      e.WebDevice.EngineVersion,
		"os_name":             url.QueryEscape(e.WebDevice.OsName),
		"os_version":          e.WebDevice.OsVersion,
		"cpu_core_num":        e.WebDevice.CpuCoreNum,
		"device_memory":       e.WebDevice.DeviceMemory,
		"platform":            e.WebDevice.Platform,
		"downlink":            e.WebDevice.Downlink,
		"effective_type":      e.WebDevice.EffectiveType,
		"round_trip_time":     e.WebDevice.RoundTripTime,
		"webid":               e.WebDevice.Webid,
		"uifid":               e.WebDevice.Uifid,
		"verify_fp":           e.WebDevice.VerifyFp,
		"fp":                  e.WebDevice.Fp,
		"msToken":             e.getMsToken(107),
	}
}

func (e *DyBaseEntity) GetParams() string {
	split := strings.Split(e.Url, "?")
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

func (e *DyBaseEntity) AppendUrlParams(name string, value interface{}) *DyBaseEntity {
	if e.Url[len(e.Url)-1] != '?' {
		e.Url += "&" + name + "=" + utils.InterfaceToString(value)
	} else {
		e.Url += name + "=" + utils.InterfaceToString(value)
	}
	return e
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
	var singUrl = viper.GetString("plugin.url")
	result, _ := http.Post(singUrl+"/dy/abogus/sign", map[string]interface{}{
		"params": params,
		"ua":     ua,
	}, "", nil, "")
	return result["aBogus"].(string)
}

func (e *DyBaseEntity) GetAcSign(Url string, acNonce string, ua string) string {
	var singUrl = viper.GetString("plugin.url")
	result, _ := http.Post(singUrl+"/dy/ac/sign", map[string]interface{}{
		"Url":     Url,
		"acNonce": acNonce,
		"ua":      ua,
	}, "", nil, "")
	return result["acSignature"].(string)
}

func (e *DyBaseEntity) GetUrl() string {
	return e.Url
}

func (e *DyBaseEntity) Sign() {
	abogus := e.GetAbogus(e.GetParams(), e.WebDevice.UserAgent)
	e.AppendUrlParams("a_bogus", abogus)
}

func (r *DyBaseEntity) SetBody(params map[string]interface{}) {
	r.Body = params
}

type DyRequest[E service.Entity] struct {
	service.Request[E]
}

func (r *DyRequest[E]) DoGet(e E, ip string) (map[string]interface{}, error) {
	e.Sign()
	result, err := http.Get(e.GetUrl(), e.GetCookieString(), e.GetHeaders(), ip)
	return result, err
}

func (r *DyRequest[E]) DoPost(e E, contentType string, ip string) (map[string]interface{}, error) {
	e.Sign()
	if contentType != "" && contentType == "application/x-www-form-urlencoded; charset=UTF-8" {
		return http.PostForm(e.GetUrl(), e.GetBody(), e.GetCookieString(), e.GetHeaders(), ip)
	}
	return http.Post(e.GetUrl(), e.GetBody(), e.GetCookieString(), e.GetHeaders(), ip)
}

func DoGet(e service.Entity) (map[string]interface{}, error) {
	requestInstance := &DyRequest[service.Entity]{}
	return requestInstance.DoGet(e, e.GetIp())
}

func DoPost(e service.Entity, params map[string]interface{}, contentType string) (map[string]interface{}, error) {
	e.SetBody(params)
	requestInstance := &DyRequest[service.Entity]{}
	return requestInstance.DoPost(e, contentType, e.GetIp())
}
