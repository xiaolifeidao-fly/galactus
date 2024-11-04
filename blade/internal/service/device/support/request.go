package support

type Entity interface {
	GetUrl() string
	GetCookieString() string
	GetHeaders() map[string]string
	GetBody() map[string]interface{}
	GetMethod() string
}

type Request[R Entity] interface {
	DoRequest(r *R) (map[string]interface{}, error)
}
