package systemuser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"template/log"
	"template/src/database"
	"time"

	"github.com/xiangjiaoflb/httpframe"
	"github.com/xiangjiaoflb/jsonlog"
)

//获取用户名密码
func getUserInfo(ctx *httpframe.Context) (user User, err error) {
	//根据数据样式解析
	switch ctx.R.Header.Get("Content-Type") {
	case "application/json":
		buf, err := ioutil.ReadAll(ctx.R.Body)
		if err != nil {
			return user, err
		}
		err = json.Unmarshal(buf, &user)
		if err != nil {
			return user, err
		}
		return user, nil
	default:
		user.Username = ctx.R.FormValue("username")
		user.Password = ctx.R.FormValue("password")
		return
	}
}

//Login 登录
func Login(ctx *httpframe.Context) {
	flog := jsonlog.Info(log.RequestLog)
	defer flog.Msg("")

	//获取用户名和密码
	user, err := getUserInfo(ctx)
	jsonlog.RequestLog(ctx.R, flog)
	if err != nil {
		jsonlog.SendJSON(flog, ctx.W, err, nil, 401)
		return
	}

	if user.Username == "" || user.Password == "" {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("参数错误"), nil, 401)
		return
	}

	//查询内存或数据库的信息
	quser, err := queryUserInfo(user.Username)
	if err != nil {
		jsonlog.SendJSON(flog, ctx.W, err, nil, 401)
		return
	}

	if quser.Username == "" || quser.Password == "" {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("系统内部错误"), nil, 401)
		return
	}

	//判断是否正确
	if !checkoutPasswork(user.Password, quser.Salt, quser.Password) {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("用户名或者密码错误"), nil, 401)
		return
	}

	//创建jwt

}

//查询内存或者数据库的用户信息
func queryUserInfo(username string) (user User, err error) {
	//判断验证的类型
	switch authtype {
	case AuthOne:
		//读数据库
		err = database.DB.First(&user, User{Username: username}).Error
		return
	case AuthTwo, AuthThree:
		//先读内存，没有再读数据库
		if v, ok := memoryuser.Load(username); ok {
			user, _ = v.(User)
			return
		}

		//读数据库
		err = database.DB.First(&user, User{Username: username}).Error
		if err != nil {
			return
		}

		//增加session
		if authtype == AuthThree {
			user.session = getSalt()
		}

		err = memoryuser.Store(username, user, time.Minute*30)
		return
	default:
		err = fmt.Errorf("系统内部错误")
		return
	}
}
