package systemuser

import (
	"fmt"
	"template/log"

	"github.com/xiangjiaoflb/httpframe"
	"github.com/xiangjiaoflb/jsonlog"
)

//LoginVerify 判断是否登陆,验证 的中间件
//读取用户token
func LoginVerify(ctx *httpframe.Context) {
	var err error
	defer func() {
		if err != nil {
			flog := jsonlog.Info(log.RequestLog)
			defer flog.Msg("")
			jsonlog.RequestLog(ctx.R, flog)
			jsonlog.SendJSON(flog, ctx.W, err, nil, 401)
		}
	}()

	//判断验证类型
	//确定签名
	var user User
	var jwtmap map[string]interface{}
	//解析 token  //验证了密码
	err = ParseJWT(gettoken(ctx), func(mapint map[string]interface{}) (interface{}, error) {
		jwtmap = mapint

		inter, ok := jwtmap[keyUsername]
		if !ok {
			return nil, fmt.Errorf("token失效")
		}
		username, ok := inter.(string)
		if !ok {
			return nil, fmt.Errorf("token失效")
		}

		if username == "" {
			return nil, fmt.Errorf("token失效")
		}

		signature := user.getuserAndSignature(username)
		if signature == "" {
			return nil, fmt.Errorf("系统内部错误")
		}
		return []byte(signature), nil
	})
	if err != nil {
		return
	}

	//验证信息 ip
	inter, ok := jwtmap[keyIP]
	if !ok {
		err = fmt.Errorf("token失效")
		return
	}
	ip, ok := inter.(string)
	if !ok {
		err = fmt.Errorf("token失效")
		return
	}

	if ip != getIP(ctx.R.RemoteAddr) {
		err = fmt.Errorf("ip错误")
		return
	}

	//检验信息 session
	if authtype == AuthThree {
		inter, ok = jwtmap[keySession]
		if !ok {
			err = fmt.Errorf("token失效")
			return
		}
		session, ok := inter.(string)
		if !ok {
			err = fmt.Errorf("token失效")
			return
		}

		if session != user.session {
			err = fmt.Errorf("账号在其他地方登录")
			return
		}
	}

	ctx.Data = &user
	ctx.Next()
}

//获取 token
func gettoken(ctx *httpframe.Context) (token string) {
	token = ctx.R.Header.Get("Authorization")
	return
}

func (user *User) getuserAndSignature(username string) (signature string) {
	//确定签名
	switch authtype {
	case AuthOne:
		signature = authstr
	case AuthTwo:
		if user.Password == "" {
			//查询数据库
			newuser, _ := queryUserInfo(username, false)
			*user = *newuser
		}
		signature = user.Password
	case AuthThree:
		if user.Password == "" {
			//查询数据库
			newuser, _ := queryUserInfo(username, false)
			*user = *newuser
		}
		signature = user.Password
	}
	return
}
