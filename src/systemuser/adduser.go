package systemuser

import (
	"encoding/json"
	"fmt"
	"template/log"
	"template/src/database"

	"github.com/xiangjiaoflb/httpframe"
	"github.com/xiangjiaoflb/jsonlog"
)

//Register 注册
func Register(ctx *httpframe.Context) {
	flog := jsonlog.Info(log.RequestLog)
	defer flog.Msg("")

	//记录请求
	bodybuf := jsonlog.RequestLog(ctx.R, flog)

	//此处可判断权限
	//puser := getCtxData(ctx)

	//判断请求方式
	if ctx.R.Method != "POST" {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("请求方式错误"), nil, 400)
		return
	}

	var user User
	err := json.Unmarshal(bodybuf, &user)
	if err != nil {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("数据错误"), nil, 400)
		return
	}

	//判断数据
	if user.Username == "" || user.Password == "" {
		jsonlog.SendJSON(flog, ctx.W, fmt.Errorf("数据错误"), nil, 400)
		return
	}

	//加密密码 存数据库
	err = RegisterUser(user)
	if err != nil {
		jsonlog.SendJSON(flog, ctx.W, err, nil, 400)
		return
	}

	jsonlog.SendJSON(flog, ctx.W, nil, nil, 200)
}

//RegisterUser 注册用户
func RegisterUser(user User) error {
	//加密密码 存数据库
	user.Password, user.Salt = encryptPasswork(user.Password)
	err := database.DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
