package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2
var MySecret = []byte("aaa");

func GenToken(username string) (string, error){
	c := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:"my-project",  // 签发人
		},
	}
   token := jwt.NewWithClaims(jwt.SigningMethodES256, c)
   if len(MySecret) == 0 {
   		return "", errors.New("token_key为空")
   }
   return token.SignedString(MySecret)
}
