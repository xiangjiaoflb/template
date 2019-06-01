package systemuser

import (
	"fmt"
	"os"
	"template/log"
	"template/src/database"
	"template/utils"
	"time"

	"github.com/xiangjiaoflb/jsonlog"
)

// AuthType 采用认证的类型
type AuthType int

//
const (
	AuthOne AuthType = iota
	AuthTwo
	AuthThree

	//第一种模式验证的签名
	authstr = "authone"

	//username key
	keyUsername = "username"

	//password key
	keyPassword = "password"

	//session key
	keySession = "session"

	//ip key
	keyIP = "ip"

	//第二三模式时刻控制最大登录量
	loginsize = 100
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

	//权限等
}

//TableName 表名
func (User) TableName() string {
	return "user"
}

func init() {
	go func() {
		for {
			time.Sleep(time.Second)
			if database.DB != nil {
				//创建表
				err := database.DB.AutoMigrate(&User{}).Error
				if err != nil {
					jsonlog.Error(log.RunLog).Err(err).Msg("")
					os.Exit(-1)
				}

				user := User{
					Username: "admin",
					Password: "admin",
				}

				//创建数据
				//加密密码 存数据库
				user.Password, user.Salt = encryptPasswork(user.Password)
				err = database.DB.FirstOrCreate(&user, &User{Username: "admin"}).Error
				if err != nil {
					jsonlog.Error(log.RunLog).Err(err).Msg("")
					os.Exit(-1)
				}
				break
			}
		}
	}()
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
			memoryuser = utils.NewMemory(loginsize)
		}
	case AuthThree:
		if memoryuser == nil {
			memoryuser = utils.NewMemory(loginsize)
		}
	default:
		return fmt.Errorf("不支持的类型")
	}

	authtype = at
	return nil
}
