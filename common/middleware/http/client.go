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

var request *gorequest.SuperAgent

func init() {
	request = InitHttpRequest()
}

func InitHttpRequest() *gorequest.SuperAgent {
	return gorequest.New().Timeout(90 * time.Second)
}

func Get(requestUrl string, cookie string, headers map[string]string, ip string) (map[string]interface{}, error) {
	req := request.Get(requestUrl)
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
	formData := url.Values{}
	for key, value := range requestBody {
		formData.Add(key, utils.InterfaceToString(value))
	}
	req := request.Post(requestUrl).Type("form").Send(formData.Encode())
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
	req := request.Post(requestUrl).Send(requestBody)
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
	req := request.Get(requestUrl)
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
