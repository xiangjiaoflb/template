package systemuser

import (
	"fmt"
	"template/log"

	"github.com/xiangjiaoflb/httpframe"
	"github.com/xiangjiaoflb/jsonlog"
)

//
func getCtxData(ctx *httpframe.Context) *User {
	puser, _ := ctx.Data.(*User)
	return puser
}

//KeepAlive 保活
func KeepAlive(ctx *httpframe.Context) {
	flog := jsonlog.Info(log.RequestLog)
	defer flog.Msg("")

	puser := getCtxData(ctx)

	if puser.Username == "" || puser.Password == "" || puser.Salt == "" {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("系统内部错误"), nil, 401)
		return
	}

	signature := puser.getuserAndSignature(puser.Username)

	//创建jwt
	token, err := CreateJWT(map[string]interface{}{
		keyUsername: puser.Username,
		keySession:  puser.session,
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
