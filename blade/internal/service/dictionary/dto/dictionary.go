package dto

import (
	"galactus/common/base/dto"
)

type DictionaryDTO struct {
	dto.BaseDTO
	Code        string `json:"code" description:"编码"`
	Value       string `json:"value" description:"值"`
	Description string `json:"description" description:"描述"`
	Type        string `json:"type" description:"类型"`
}
