package repository

import (
	"galactus/common/middleware/db"
)

type Dictionary struct {
	db.BaseEntity
	Code        string `orm:"column(code);size(256);null" description:"编码"`
	Value       string `orm:"column(value);size(256);null" description:"值"`
	Description string `orm:"column(description);size(256);null" description:"描述"`
	Type        string `orm:"column(type);size(256);null" description:"类型"`
}

func (d *Dictionary) TableName() string {
	return "dictionary"
}

type DictionaryRepository struct {
	db.Repository[*Dictionary]
}

func (r *DictionaryRepository) GetByCode(code string) (*Dictionary, error) {
	dict, err := r.GetOne("select * from dictionary where code = ? and active = 1", code)
	return dict, err
}

func (r *DictionaryRepository) GetByType(typeStr string) ([]*Dictionary, error) {
	dicts, err := r.GetList("select * from dictionary where type = ? and active = 1", typeStr)
	return dicts, err
}
