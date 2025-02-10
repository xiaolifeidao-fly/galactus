package ip

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"galactus/blade/internal/consts"
	"galactus/blade/internal/service/dictionary"
)

// ProxyIP 代理IP信息
type ProxyIP struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Address  string `json:"adr"`
	Timeout  int    `json:"timeout"`
	Cometime int    `json:"cometime"`
}

type ZDYHttpProxyService struct {
	DictionaryService dictionary.DictionaryService
}

var (
	defaultZDYInstance *ZDYHttpProxyService
	zdyOnce            sync.Once
)

func GetDefaultZDYHttpProxyService() *ZDYHttpProxyService {
	zdyOnce.Do(func() {
		defaultZDYInstance = &ZDYHttpProxyService{
			DictionaryService: dictionary.NewDictionaryService(),
		}
	})
	return defaultZDYInstance
}

func (s *ZDYHttpProxyService) GetIp() string {
	ip, err := s.GetUserIpByProxyType(consts.SceneAuditLike)
	if err != nil {
		return ""
	}
	if len(ip) == 0 {
		return ""
	}
	return ip[0].IP + ":" + strconv.Itoa(ip[0].Port)
}

func (s *ZDYHttpProxyService) GetUserIpByProxyType(scene consts.Scene) ([]ProxyIP, error) {
	var url string
	proxyUrl, err := s.DictionaryService.GetByCode(scene.GetProxyRequestUrl())
	if err != nil {
		return nil, err
	}
	proxySuffix, err := s.DictionaryService.GetByCode(scene.GetProxyRequestSuffix())
	if err != nil {
		return nil, err
	}
	proxyApi, err := s.DictionaryService.GetByCode(scene.GetProxyRequestApi())
	if err != nil {
		return nil, err
	}
	proxyAkey, err := s.DictionaryService.GetByCode(scene.GetProxyRequestAkey())
	if err != nil {
		return nil, err
	}
	url = fmt.Sprintf(proxyUrl.Value+proxySuffix.Value,
		proxyApi.Value,
		proxyAkey.Value)
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
