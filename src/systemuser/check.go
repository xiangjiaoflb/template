package systemuser

import (
	"template/log"

	"github.com/xiangjiaoflb/httpframe"
	"github.com/xiangjiaoflb/jsonlog"
)

//LoginVerify 判断是否登陆,校验用户登陆的中间件
//读取用户token
//判断token值（ip 用户名）
func LoginVerify(ctx *httpframe.Context) {
	flog := jsonlog.Info(log.RequestLog)
	defer flog.Msg("")

	//jsonlog.RequestLog(ctx.R, flog)

	//获取 token

	//解析 token

	//验证信息 ip username

	ctx.Next()
}

//获取 token
func gettoken(){
	
}

//解析 token

//验证信息 ip username
