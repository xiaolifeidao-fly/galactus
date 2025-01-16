package pb

import (
	"SignIos/utils"
	"encoding/json"
)

// XArgusToJsonString 转成字符串
func XArgusToJsonString(info XArgus) (string, error) {
	res, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func XGorgonToJsonString(info utils.XGorGonInfo) (string, error) {
	res, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func ToJsonString(info interface{}) (string, error) {
	res, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
