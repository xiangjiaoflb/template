package systemuser

import (
	"fmt"
	"template/utils"
)

// AuthType 采用认证的类型
type AuthType int

//
const (
	AuthOne AuthType = iota
	AuthTwo
	AuthThree

	//第一种验证的签名
	authstr = "authone"
)

var (
	authtype = AuthOne

	//用户内存
	memoryuser *utils.Memory
)

// User ...
type User struct {
	Username string `gorm:"column:username;not null;unique" json:"username"` //用户名
	Password string `gorm:"column:password;not null" json:"password"`        //密码
	Salt     string `gorm:"column:salt;not null" json:"salt"`                //盐

	session string
}

//TableName 表名
func (User) TableName() string {
	return "user"
}

//SetAuthType 设置验证方式
func SetAuthType(at AuthType) error {
	switch at {
	case AuthOne:
		authtype = at
		if memoryuser != nil {
			memoryuser.Close()
		}
		memoryuser = nil
	case AuthTwo:
		if memoryuser == nil {
			memoryuser = utils.NewMemory(1000)
		}
	case AuthThree:
		if memoryuser == nil {
			memoryuser = utils.NewMemory(1000)
		}
	default:
		return fmt.Errorf("不支持的类型")
	}

	authtype = at
	return nil
}