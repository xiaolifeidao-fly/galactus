package http

import (
	"encoding/json"
	"galactus/common/utils"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

type SuperAgentPool struct {
	pool chan *gorequest.SuperAgent
}

func NewSuperAgentPool(size int) *SuperAgentPool {
	p := &SuperAgentPool{
		pool: make(chan *gorequest.SuperAgent, size),
	}
	for i := 0; i < size; i++ {
		p.pool <- gorequest.New().Timeout(90 * time.Second)
	}
	return p
}

func (p *SuperAgentPool) Get() *gorequest.SuperAgent {
	return <-p.pool
}

func (p *SuperAgentPool) Put(req *gorequest.SuperAgent) {
	req.ClearSuperAgent()
	req.Timeout(90 * time.Second)
	p.pool <- req
}

var superAgentPool = NewSuperAgentPool(10) // 设置池的大小

func Get(requestUrl string, cookie string, headers map[string]string, ip string) (map[string]interface{}, error) {
	req := superAgentPool.Get()
	defer superAgentPool.Put(req)

	req.Get(requestUrl)
	if ip != "" {
		proxyUrl, err := url.Parse(ip)
		if err == nil {
			req.Proxy(proxyUrl.String())
		}
	}
	if cookie != "" {
		req.Set("cookie", cookie)
	}
	for key, value := range headers {
		req.Set(key, value)
	}
	_, body, errs := req.End()
	if len(errs) > 0 {
		log.Fatalf("请求失败: %v", errs)
	}
	result := map[string]interface{}{}
	json.Unmarshal([]byte(body), &result)
	return result, nil
}

func PostForm(requestUrl string, requestBody map[string]interface{}, cookie string, headers map[string]string, ip string) (map[string]interface{}, error) {
	req := superAgentPool.Get()
	defer superAgentPool.Put(req)

	formData := url.Values{}
	for key, value := range requestBody {
		formData.Add(key, utils.InterfaceToString(value))
	}
	req.Post(requestUrl).Type("form").Send(formData.Encode())
	if ip != "" {
		proxyUrl, err := url.Parse(ip)
		if err == nil {
			req.Proxy(proxyUrl.String())
		}
	}
	if cookie != "" {
		req.Set("Cookie", cookie)
	}
	for key, value := range headers {
		req.Set(key, value)
	}
	_, body, errs := req.End()
	if len(errs) > 0 {
		log.Fatalf("请求失败: %v", errs)
	}
	result := map[string]interface{}{}
	json.Unmarshal([]byte(body), &result)
	return result, nil
}

func Post(requestUrl string, requestBody map[string]interface{}, cookie string, headers map[string]string, ip string) (map[string]interface{}, error) {
	req := superAgentPool.Get()
	defer superAgentPool.Put(req)

	req.Post(requestUrl).Send(requestBody)
	if ip != "" {
		proxyUrl, err := url.Parse(ip)
		if err == nil {
			req.Proxy(proxyUrl.String())
		}
	}
	if cookie != "" {
		req.Set("Cookie", cookie)
	}
	for key, value := range headers {
		req.Set(key, value)
	}
	_, body, errs := req.End()
	if len(errs) > 0 {
		log.Fatalf("请求失败: %v", errs)
	}
	result := map[string]interface{}{}
	json.Unmarshal([]byte(body), &result)
	return result, nil
}

func GetToResponse(requestUrl string, cookie string, headers map[string]string, ip string) (*http.Response, error) {
	req := superAgentPool.Get()
	defer superAgentPool.Put(req)

	req.Get(requestUrl)
	if ip != "" {
		proxyUrl, err := url.Parse(ip)
		if err == nil {
			req.Proxy(proxyUrl.String())
		}
	}
	if cookie != "" {
		req.Set("Cookie", cookie)
	}
	for key, value := range headers {
		req.Set(key, value)
	}
	resp, _, errs := req.End()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return resp, nil
}
