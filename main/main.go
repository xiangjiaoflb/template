package main

import (
	"net/http"
	"os"
	"template/log"
	"template/route"
	"template/src/database"
	"template/src/systemuser"

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

	//注册路由 不走中间件
	httpframe.RegisterHandle(servermux, nil, &route.Api{})

	//注册走中间件的路由
	httpframe.RegisterHandle(servermux, []httpframe.HandlerFunc{systemuser.LoginVerify}, &route.User{})

	//其他路由走中间件
	servermux.HandleFunc("/", httpframe.NewMiddleware(append([]httpframe.HandlerFunc{systemuser.LoginVerify},
		func(ctx *httpframe.Context) { http.FileServer(http.Dir(".")).ServeHTTP(ctx.W, ctx.R) })).HandleFunc)

	err := http.ListenAndServe(":8888", servermux)
	if err != nil {
		jsonlog.Error(log.RunLog).Err(err).Msg("")
	}
}

//创建数据库连接和表
func init() {
	db, err := database.Open("root:root@tcp(192.168.216.129:3306)/mydata?charset=utf8&parseTime=True")
	if err != nil {
		jsonlog.Error(log.RunLog).Err(err).Msg("")
		os.Exit(-1)
	}
	if DEBUG {
		db.LogMode(true)
	}

	//创建表
	err = db.AutoMigrate(&systemuser.User{}).Error
	if err != nil {
		jsonlog.Error(log.RunLog).Err(err).Msg("")
		os.Exit(-1)
	}

	user := systemuser.User{
		Username: "admin",
		Password: "admin",
	}
	//创建数据
	systemuser.RegisterUser(user)
}

//
var (
	DEBUG = true
)
