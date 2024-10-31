package db

import (
	"gorm.io/gorm"
)

// NewRepository 创建新仓库实例
func NewRepository[R any]() *R {
	var repoValue *R = new(R)
	if repo, ok := any(repoValue).(interface{ SetDb(*gorm.DB) }); ok {
		repo.SetDb(Db)
	}
	return repoValue
}
