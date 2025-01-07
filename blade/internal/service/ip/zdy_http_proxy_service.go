package ip

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ZDYHttpProxyService struct {
	url     string
	key     string
	api     string
	ipQueue *IPQueue
}

type IPQueue struct {
	items []map[string]interface{}
	mu    sync.Mutex
}

func NewIPQueue() *IPQueue {
	return &IPQueue{
		items: make([]map[string]interface{}, 0),
	}
}

func (q *IPQueue) Enqueue(item map[string]interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *IPQueue) Dequeue() (map[string]interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

var (
	defaultZDYInstance *ZDYHttpProxyService
	zdyOnce            sync.Once
)

func InitDefaultZDYHttpProxyService(url, key, api string) error {
	zdyOnce.Do(func() {
		defaultZDYInstance = &ZDYHttpProxyService{
			url:     url,
			key:     key,
			api:     api,
			ipQueue: NewIPQueue(),
		}
	})
	return nil
}

func GetDefaultZDYHttpProxyService() *ZDYHttpProxyService {
	if defaultZDYInstance == nil {
		panic("ZDYHttpProxyService is not initialized. Call InitDefaultZDYHttpProxyService first.")
	}
	return defaultZDYInstance
}

func (s *ZDYHttpProxyService) GetUserIpByProxyType(fetchNum int) (map[string]interface{}, error) {
	if len(s.ipQueue.items) == 0 {
		err := s.init(fetchNum)
		if err != nil {
			return nil, err
		}
	}
	ip, ok := s.ipQueue.Dequeue()
	if !ok {
		return nil, errors.New("failed to get IP from queue")
	}
	return ip, nil
}

func (s *ZDYHttpProxyService) init(fetchNum int) error {
	url := fmt.Sprintf("%s?api=%s&akey=%s&pro=1&order=1&type=3&count=%d", s.url, s.api, s.key, fetchNum)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if code, ok := result["code"].(string); !ok || code != "10001" {
		return errors.New("failed to fetch IPs")
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return errors.New("invalid data format")
	}

	proxyList, ok := data["proxy_list"].([]interface{})
	if !ok {
		return errors.New("invalid proxy list format")
	}

	for _, item := range proxyList {
		ipData, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		expireTime, err := time.Parse(time.RFC3339, ipData["expireTime"].(string))
		if err != nil {
			return err
		}
		ipData["expireTime"] = expireTime
		s.ipQueue.Enqueue(ipData)
	}
	return nil
}
