package main

import (
	"fmt"
	"net/http"
	"template/log"
	"template/route"

	"github.com/xiangjiaoflb/httpframe"
	"github.com/xiangjiaoflb/jsonlog"
)

//版本号等
var (
	// VERSION 版本
	VERSION = "1.0.1"

	// BUILDTIME 编译时间
	BUILDTIME = ""

	// GOVERSION go 版本
	GOVERSION = ""

	// GITHASH 代码hash值
	GITHASH = ""
)

func main() {
	//程序启动打印日志
	jsonlog.Info(log.RunLog).
		Str(VERSION, VERSION).
		Str(BUILDTIME, BUILDTIME).
		Str(GOVERSION, GOVERSION).
		Str(GITHASH, GITHASH).Msg("begin run!")

	servermux := http.NewServeMux()

	//注册路由
	httpframe.RegisterHandle(servermux, nil, &route.User{})

	//其他路由走中间件
	servermux.HandleFunc("/", httpframe.NewMiddleware(append([]httpframe.HandlerFunc{},
		func(ctx *httpframe.Context) { http.FileServer(http.Dir(".")).ServeHTTP(ctx.W, ctx.R) })).HandleFunc)

	err := http.ListenAndServe(":8888", servermux)
	if err != nil {
		fmt.Println(err)
	}
}
