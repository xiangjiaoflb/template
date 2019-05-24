package systemuser

import (
	"template/log"

	"github.com/xiangjiaoflb/httpframe"
	"github.com/xiangjiaoflb/jsonlog"
)

//LogOut 退出登录
func LogOut(ctx *httpframe.Context) {
	flog := jsonlog.Info(log.RequestLog)
	defer flog.Msg("")
	jsonlog.RequestLog(ctx.R, flog)

	//删除内存存的信息
	

	jsonlog.SendJSON(flog, ctx.W, nil, nil, 200)
}
