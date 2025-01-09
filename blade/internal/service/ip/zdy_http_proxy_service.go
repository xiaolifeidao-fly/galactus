package ip

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"galactus/common/middleware/vipper"
)

// ProxyIP 代理IP信息
type ProxyIP struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Address  string `json:"adr"`
	Timeout  int    `json:"timeout"`
	Cometime int    `json:"cometime"`
}

type ZDYHttpProxyService struct{}

var (
	defaultZDYInstance *ZDYHttpProxyService
	zdyOnce            sync.Once
)

func GetDefaultZDYHttpProxyService() *ZDYHttpProxyService {
	zdyOnce.Do(func() {
		defaultZDYInstance = &ZDYHttpProxyService{}
	})
	return defaultZDYInstance
}

func (s *ZDYHttpProxyService) GetUserIpByProxyType(fetchNum int) ([]ProxyIP, error) {
	var url string
	if fetchNum > 0 {
		url = fmt.Sprintf("%s?api=%s&akey=%s&pro=1&order=1&type=3&count=%d",
			vipper.GetString("zdy.url"),
			vipper.GetString("zdy.api"),
			vipper.GetString("zdy.key"),
			fetchNum)
	} else {
		url = fmt.Sprintf("%s?api=%s&akey=%s&pro=1&order=1&type=3",
			vipper.GetString("zdy.url"),
			vipper.GetString("zdy.api"),
			vipper.GetString("zdy.key"))
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Count     int       `json:"count"`
			ProxyList []ProxyIP `json:"proxy_list"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != "10001" {
		return nil, errors.New("failed to fetch IPs")
	}

	if len(result.Data.ProxyList) == 0 {
		return nil, errors.New("empty proxy list")
	}

	return result.Data.ProxyList, nil
}
