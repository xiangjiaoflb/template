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
	quser, err := queryUserInfo(user.Username, true)
	if err != nil {
		jsonlog.SendJSON(flog, ctx.W, err, nil, 401)
		return
	}

	if quser.Username == "" || quser.Password == "" || quser.Salt == "" {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("系统内部错误"), nil, 401)
		return
	}

	//判断是否正确
	if !checkoutPasswork(user.Password, quser.Salt, quser.Password) {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("用户名或者密码错误"), nil, 401)
		return
	}

	//确定签名
	signature := quser.getuserAndSignature(quser.Username)

	//创建jwt
	token, err := CreateJWT(map[string]interface{}{
		keyUsername: quser.Username,
		keySession:  quser.session,
		keyIP:       getIP(ctx.R.RemoteAddr),
	}, signature)
	if err != nil {
		jsonlog.SendJSON(flog, ctx.W, err, nil, 401)
		return
	}

	jsonlog.SendJSON(flog, ctx.W, nil, map[string]interface{}{
		"token": token,
	}, 200)
}

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
		user.Username = ctx.R.FormValue(keyUsername)
		user.Password = ctx.R.FormValue(keyPassword)
		return
	}
}

//查询内存或者数据库的用户信息
func queryUserInfo(username string, replacesession bool) (puser *User, err error) {
	puser = new(User)

	//判断验证的类型
	switch authtype {
	case AuthOne:
		//读数据库
		err = database.DB.First(puser, User{Username: username}).Error
		return
	case AuthTwo, AuthThree:
		//先读内存，没有再读数据库
		if v, ok := memoryuser.Load(username); ok {
			puser, _ = v.(*User)
			//增加session
			if replacesession && authtype == AuthThree {
				puser.session = getSalt()
			}
			return
		}

		//读数据库
		err = database.DB.First(puser, &User{Username: username}).Error
		if err != nil {
			return
		}

		//增加session //保证模式2切换到3时之前的token也能用
		if replacesession && authtype == AuthThree {
			puser.session = getSalt()
		}

		err = memoryuser.Store(username, puser, time.Minute*30)
		return
	default:
		err = fmt.Errorf("系统内部错误")
		return
	}
}
