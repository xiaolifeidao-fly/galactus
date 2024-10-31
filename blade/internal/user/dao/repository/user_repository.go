package repository

import (
	"galactus/blade/internal/user/dao/po"
	"galactus/common/middleware/db"
)

type UserRepository struct {
	db.Repository[*po.User]
}
