package route

import (
	"template/src/systemuser"

	"github.com/xiangjiaoflb/httpframe"
)

// Api ...
type Api struct{}

//Login 登录
func (Api) Login(ctx *httpframe.Context) {
	systemuser.Login(ctx)
}

//LoginOut 退出登录
func (Api) LoginOut(ctx *httpframe.Context) {
}

// User ...
type User struct{}

//Register 注册账号
func (User) Register(ctx *httpframe.Context) {
	systemuser.Register(ctx)
}
