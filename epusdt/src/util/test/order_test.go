package test

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"testing"
)

// 创建交易请求结构体
type CreateTransactionRequest struct {
	OrderId     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	NotifyUrl   string  `json:"notify_url"`
	RedirectUrl string  `json:"redirect_url"`
	Signature   string  `json:"signature"`
}

// 生成签名
func generateSignature(params map[string]string, token string) string {
	// 1. 获取所有键并排序
	var keys []string
	for k := range params {
		if k != "signature" && params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 2. 按排序后的键拼接字符串
	var signParts []string
	for _, k := range keys {
		signParts = append(signParts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signStr := strings.Join(signParts, "&")

	// 3. 拼接token并计算MD5
	signStr = signStr + token
	hash := md5.Sum([]byte(signStr))
	return fmt.Sprintf("%x", hash)
}

func Test_request_create_transaction(t *testing.T) {
	// API配置
	apiEndpoint := "http://127.0.0.1:8000/api/v1/order/create-transaction"
	token := "api_auth_token_roc" // 在.env文件中配置的api接口认证token

	// 准备请求参数
	params := map[string]string{
		"order_id":     "20220201030210321",
		"amount":       "42",
		"notify_url":   "http://example.com/notify",
		"redirect_url": "http://example.com/redirect",
	}

	// 生成签名
	signature := generateSignature(params, token)

	// 构建请求数据
	reqData := CreateTransactionRequest{
		OrderId:     params["order_id"],
		Amount:      42.00,
		NotifyUrl:   params["notify_url"],
		RedirectUrl: params["redirect_url"],
		Signature:   signature,
	}

	// 转换为JSON
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("读取响应失败: %v", err)
	}

	// 打印响应结果
	t.Logf("响应状态码: %d", resp.StatusCode)
	t.Logf("响应内容: %s", string(body))

	// 如果需要，可以解析响应JSON
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("解析响应JSON失败: %v", err)
	}

	// 验证响应中的关键字段
	if statusCode, ok := response["status_code"].(float64); !ok || statusCode != 200 {
		t.Errorf("请求失败，状态码: %v", response["status_code"])
	}

	// 打印成功创建的订单信息
	if data, ok := response["data"].(map[string]interface{}); ok {
		t.Logf("交易ID: %v", data["trade_id"])
		t.Logf("支付地址: %v", data["payment_url"])
		t.Logf("实际支付金额(USDT): %v", data["actual_amount"])
	}
}

//order_test.go:103: 响应状态码: 200
//order_test.go:104: 响应内容: {"status_code":200,"message":"success","data":{"trade_id":"202412211734774167045146","order_id":"20220201030210321","amount":42,"actual_amount":5.76,"token":"TMe5Ev1DfNk69eCLeyWHohECAXJY5v62a1","expiration_time":1734774768,"payment_url":"https://dujiaoka.com/pay/checkout-counter/202412211734774167045146"},"request_id":"ff30cd03-b1e2-47ea-bd77-878a6e1407da"}
//order_test.go:119: 交易ID: 202412211734774167045146
//order_test.go:120: 支付地址: https://dujiaoka.com/pay/checkout-counter/202412211734774167045146
//order_test.go:121: 实际支付金额(USDT): 5.76
