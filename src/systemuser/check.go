package systemuser

import (
	"fmt"
	"template/log"

	"github.com/xiangjiaoflb/httpframe"
	"github.com/xiangjiaoflb/jsonlog"
)

func getusername(ctx *httpframe.Context) string {
	return ctx.R.FormValue(keyUsername)
}

//LoginVerify 判断是否登陆,验证 的中间件
//读取用户token
func LoginVerify(ctx *httpframe.Context) {
	flog := jsonlog.Info(log.RequestLog)
	defer flog.Msg("")

	//判断验证类型
	//确定签名
	var user User
	signature := user.getuserAndSignature(getusername(ctx))

	//解析 token  //验证了密码
	jwtmap, err := ParseJWT(gettoken(ctx), signature)
	if err != nil {
		jsonlog.SendJSON(flog, ctx.W, err, nil, 401)
		return
	}

	//验证信息 ip
	inter, ok := jwtmap[keyIP]
	if !ok {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("token失效"), nil, 401)
		return
	}
	ip, ok := inter.(string)
	if !ok {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("token失效"), nil, 401)
		return
	}

	if ip != getIP(ctx.R.RemoteAddr) {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("ip错误"), nil, 401)
		return
	}

	//检验信息 session
	if authtype == AuthThree {
		inter, ok = jwtmap[keySession]
		if !ok {
			jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("token失效"), nil, 401)
			return
		}
		session, ok := inter.(string)
		if !ok {
			jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("token失效"), nil, 401)
			return
		}

		if session != user.session {
			jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("账号在其他地方登录"), nil, 401)
			return
		}
	}

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
			newuser, _ := queryUserInfo(username)
			*user = newuser
		}
		signature = user.Password
	case AuthThree:
		if user.Password == "" {
			//查询数据库
			newuser, _ := queryUserInfo(username)
			*user = newuser
		}
		signature = user.Password
	}
	return
}
