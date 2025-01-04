package repository

import (
	"galactus/common/middleware/db"
	"time"
)

type ProxyIp struct {
	db.BaseEntity
	Type       string    `orm:"column(type);size(255);null" description:"IP类型"`
	Ip         string    `orm:"column(ip);size(255);null" description:"IP地址"`
	ExpireTime time.Time `orm:"column(expire_time);type(datetime);null" description:"过期时间"`
	ApiName    string    `orm:"column(api_name);size(255);null" description:"API名称"`
	ApiKey     string    `orm:"column(api_key);size(255);null" description:"API密钥"`
}

func (p *ProxyIp) TableName() string {
	return "proxy_ip"
}

type ProxyIpRepository struct {
	db.Repository[*ProxyIp]
}

func (p *ProxyIpRepository) GetByIp(ip string) (*ProxyIp, error) {
	proxyIp, err := p.GetOne("select * from proxy_ip where ip = ? and active = 1", ip)
	return proxyIp, err
}

func (p *ProxyIpRepository) GetByType(ipType string) ([]*ProxyIp, error) {
	proxyIps, err := p.GetList("select * from proxy_ip where type = ? and active = 1", ipType)
	return proxyIps, err
}
