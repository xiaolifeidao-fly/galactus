package test

import (
	"galactus/blade/internal/consts"
	"galactus/blade/internal/service/device/biz"
	biz2 "galactus/blade/internal/service/ip/biz"
	"galactus/blade/internal/service/ip/dto"
	"galactus/common/middleware/routers"
	"sync"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	*routers.BaseHandler
	webDeviceManager *biz.WebDeviceManager
	ipManager        *biz2.IpManager
}

func NewTestHandler() *TestHandler {
	return &TestHandler{
		webDeviceManager: biz.GetDefaultWebDeviceManager(),
		ipManager:        biz2.GetDefaultIpManager(),
	}
}

func (h *TestHandler) RegisterHandler(engine *gin.RouterGroup) {
	//用户列表
	engine.GET("/test/device", h.getDevice)
	//获取ip
	engine.GET("/test/ip", h.getIp)
	//并发获取多个ip
	engine.GET("/test/ips", h.getMultipleIps)
}

func (h *TestHandler) getDevice(context *gin.Context) {
	//device, err := h.webDeviceManager.GetWebDevice()
	device, err := biz.GetDefaultWebDeviceManager().GetWebDevice(consts.SceneAuditLike)
	routers.ToJson(context, device, err)
}

func (h *TestHandler) getIp(context *gin.Context) {
	ip, err := biz2.GetDefaultIpManager().GetIp(consts.SceneCollectDevice)
	routers.ToJson(context, ip, err)
}

func (h *TestHandler) getMultipleIps(context *gin.Context) {
	const numGoroutines = 10
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results = make([]*dto.ProxyIpDTO, 0, numGoroutines)
	)

	// 创建一个通道来收集错误
	errChan := make(chan error, numGoroutines)

	// 启动5个goroutine并发获取IP
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ip, err := biz2.GetDefaultIpManager().GetIp(consts.SceneCollectDevice)
			if err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			results = append(results, ip)
			mu.Unlock()
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误发生
	for err := range errChan {
		if err != nil {
			routers.ToJson(context, nil, err)
			return
		}
	}

	routers.ToJson(context, results, nil)
}
