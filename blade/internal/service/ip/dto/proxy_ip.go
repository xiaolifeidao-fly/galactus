package dto

import (
	"galactus/common/base/dto"
	"time"
)

type ProxyIpDTO struct {
	dto.BaseDTO
	Type       string    `json:"type" description:"IP类型"`
	Ip         string    `json:"ip" description:"IP地址"`
	ExpireTime time.Time `json:"expireTime" description:"过期时间"`
	ApiName    string    `json:"apiName" description:"API名称"`
	ApiKey     string    `json:"apiKey" description:"API密钥"`
}
