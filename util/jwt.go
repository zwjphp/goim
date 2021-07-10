package util

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"net/http"
	"time"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2
var MySecret = []byte("夏天夏天悄悄过去")

// GenToken 生成JWT
func GenToken(username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "my-project",                               // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	if len(MySecret) == 0 {
		return "", errors.New("token_key为空")
	}
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// AuthHandler: 获取Token
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	r.ParseForm()
	// 检查提供的凭据-如果将这些凭据存储在数据库中，则查询将在此处进行检查。
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if username != "myusername" || password != "mypassword" {
		io.WriteString(w, `{"error":"账号或密码错误"}`)
		return
	}
	tokenString, _ := GenToken(username)
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("令牌无效")
}

// 验证token
func JWTAuthMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	r.ParseForm()
	authHeader := r.Form.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"code":-2,"msg":"token不存在"}`)
		return
	}
	mc, err := ParseToken(authHeader)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"code":-3,"msg":"无效的token"}`)
		return
	}
	fmt.Println(mc.Username)
}














