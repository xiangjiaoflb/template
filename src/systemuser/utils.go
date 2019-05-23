package systemuser

import (
	"crypto/md5"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
)

//getSalt 获取盐
func getSalt() (salt string) {
	return xid.New().String()
}

//EncryptPasswork 加密密码
func encryptPasswork(passwork string) (newpasswork string, salt string) {
	salt = getSalt()
	newpasswork = fmt.Sprintf("%x", md5.Sum([]byte(passwork+salt)))
	return
}

//CheckoutPasswork 检查密码是否正确
func checkoutPasswork(passwork string, salt string, passworkMd5 string) bool {
	return passworkMd5 == fmt.Sprintf("%x", md5.Sum([]byte(passwork+salt)))
}

// CreateJWT 创建jwt字符串
func CreateJWT(username interface{}, signature string) (jwtstr string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	//自定义claim
	token.Claims = jwt.MapClaims{
		"username": username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Second * time.Duration(60*30)).Unix(),
	}

	return token.SignedString([]byte(signature))
}

//ParseJWT 解析jwt字符串
func ParseJWT(tokenss string, signature string) (jwtmap jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenss, func(token *jwt.Token) (interface{}, error) {
		return []byte(signature), nil
	})
	if err != nil {
		return
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = fmt.Errorf("cannot convert claim to mapclaim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = fmt.Errorf("token is invalid")
		return
	}

	jwtmap = claim
	return
}
