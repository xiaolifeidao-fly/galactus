package consts

// Scene 场景类型
type Scene int

const (
	// 场景定义
	SceneCollectDevice Scene = iota // 搜集设备
	SceneAuditLike                  // 审核点赞
	SceneCurrentValue               // 当前值
	SceneAuditFollow                // 审核关注
)

// DeviceConfigType 设备配置类型
type DeviceConfigType int

const (
	DeviceCurrentIndex DeviceConfigType = iota // 设备当前值
	DeviceMaxUse                               // 设备最大使用次数
	DeviceIdRange                              // 设备一次取的数量
	DeviceMinId                                // 设备最小ID
	DeviceMaxId                                // 设备最大ID
	ProxyRequestUrl                            // 代理ip的请求url
	ProxyRequestApi                            // 请求参数api
	ProxyRequestAkey                           // 请求参数akey
)

// 场景名称映射
var sceneNames = map[Scene]string{
	SceneCollectDevice: "COLLECT",
	SceneAuditLike:     "AUDIT_LIKE",
	SceneCurrentValue:  "CURRENT",
	SceneAuditFollow:   "AUDIT_FOLLOW",
}

// GetSceneName 获取场景名称
func (s Scene) GetSceneName() string {
	return sceneNames[s]
}

// GetDeviceCurrentIndex 获取设备当前值配置key
func (s Scene) GetDeviceCurrentIndex() string {
	return "DEVICE_" + s.GetSceneName() + "_CURRENT_INDEX"
}

// GetDeviceMaxUse 获取设备最大使用次数配置key
func (s Scene) GetDeviceMaxUse() string {
	return "DEVICE_" + s.GetSceneName() + "_MAX_USE"
}

// GetDeviceIdRange 获取设备一次取的数量配置key
func (s Scene) GetDeviceIdRange() string {
	return "DEVICE_" + s.GetSceneName() + "_ID_RANGE"
}

// GetProxyRequestUrl 获取代理IP请求URL配置key
func (s Scene) GetProxyRequestUrl() string {
	return "PROXY_" + s.GetSceneName() + "_REQUEST_URL"
}

// GetProxyRequestApi 获取代理IP请求API配置key
func (s Scene) GetProxyRequestApi() string {
	return "PROXY_" + s.GetSceneName() + "_REQUEST_API"
}

// GetProxyRequestAkey 获取代理IP请求AKEY配置key
func (s Scene) GetProxyRequestAkey() string {
	return "PROXY_" + s.GetSceneName() + "_REQUEST_AKEY"
}

// GetDeviceMinId 获取设备最小ID配置key
func (s Scene) GetDeviceMinId() string {
	return "DEVICE_" + s.GetSceneName() + "_MIN_ID"
}

// GetDeviceMaxId 获取设备最大ID配置key
func (s Scene) GetDeviceMaxId() string {
	return "DEVICE_" + s.GetSceneName() + "_MAX_ID"
}

// GetConfigByType 根据配置类型获取对应的配置key
func (s Scene) GetConfigByType(configType DeviceConfigType) string {
	switch configType {
	case DeviceCurrentIndex:
		return s.GetDeviceCurrentIndex()
	case DeviceMaxUse:
		return s.GetDeviceMaxUse()
	case DeviceIdRange:
		return s.GetDeviceIdRange()
	case DeviceMinId:
		return s.GetDeviceMinId()
	case DeviceMaxId:
		return s.GetDeviceMaxId()
	case ProxyRequestUrl:
		return s.GetProxyRequestUrl()
	case ProxyRequestApi:
		return s.GetProxyRequestApi()
	case ProxyRequestAkey:
		return s.GetProxyRequestAkey()
	default:
		return ""
	}
}

// 使用示例：
/*
func main() {
    // 获取搜集设备场景的设备当前值配置
    collectCurrentIndex := SceneCollectDevice.GetDeviceCurrentIndex()
    // 输出: DEVICE_COLLECT_CURRENT_INDEX

    // 获取审核点赞场景的设备最大使用次数配置
    auditLikeMaxUse := SceneAuditLike.GetDeviceMaxUse()
    // 输出: DEVICE_AUDIT_LIKE_MAX_USE

    // 使用统一的获取方法
    auditFollowRange := SceneAuditFollow.GetConfigByType(DeviceIdRange)
    // 输出: DEVICE_AUDIT_FOLLOW_ID_RANGE

    // 获取当前值场景的代理请求URL
    currentProxyUrl := SceneCurrentValue.GetConfigByType(ProxyRequestUrl)
    // 输出: PROXY_CURRENT_REQUEST_URL
}
*/
