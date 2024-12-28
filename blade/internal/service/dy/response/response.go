package response

var (
	DELETE       = "DELETE"
	SECRET       = "SECRET"
	UN_AUTHORIZE = "UN_AUTHORIZE"
	ERROR        = "ERROR"
	NOT_GET_DATA = "NOT_GET_DATA"
	SUCCESS      = "SUCCESS"
)

var (
	HS_UID_TYPE = "HS"
	DY_UID_TYPE = "DY"
)

type BaseItem struct {
	DataStatus string `json:"dataStatus"`
}

type ConvertItemDTO struct {
	BaseItem
	ConvertValue string                 `json:"convertValue"`
	Property     map[string]interface{} `json:"property"`
}

type ExtItemDTO struct {
	BaseItem
	BusinessId string                 `json:"businessId"`
	NowNum     int64                  `json:"nowNum"`
	Name       string                 `json:"name"`
	ExtParams  map[string]interface{} `json:"extParams"`
	Property   map[string]interface{} `json:"property"`
	Uid        string                 `json:"uid"`
}

type ConvertUrlItemDTO struct {
	BaseItem
	UidType string `json:"uidType"`
	Uid     string `json:"uid"`
}
