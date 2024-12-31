package service

type Entity interface {
	GetUrl() string
	GetCookieString() string
	GetHeaders() map[string]string
	GetBody() map[string]interface{}
	SetBody(params map[string]interface{})
	GetMethod() string
	Sign()
	GetIp() string
}

type Request[R Entity] interface {
	DoGet(r *R) (map[string]interface{}, error)
	DoPost(r *R) (map[string]interface{}, error)
}
