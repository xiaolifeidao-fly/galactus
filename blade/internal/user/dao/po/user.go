package po

import "galactus/common/middleware/db"

type User struct {
	db.BaseEntity
	Channel  string `orm:"column(channel);size(32);null" description:"渠道"`
	Mobile   string `orm:"column(mobile);size(32);null" description:"手机号"`
	Nickname string `orm:"column(nickname);size(32);null" description:"昵称"`
	Password string `orm:"column(password);size(32);null" description:"密码"`
	Sex      string `orm:"column(sex);size(32);null" description:"性别"`
	Username string `orm:"column(username);size(32);null" description:"用户名"`
}

func (e *User) TableName() string {
	return "user"
}
