package route

import (
	"template/src/systemuser"

	"github.com/xiangjiaoflb/httpframe"
)

// User ...
type User struct{}

//Login 登录
func (User) Login(ctx *httpframe.Context) {
	systemuser.Login(ctx)
}

//LoginOut 退出登录
func (User) LoginOut(ctx *httpframe.Context) {
}
