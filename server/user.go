package server

import (
	"errors"
	"fmt"
	"goim/model"
	"goim/util"
	"github.com/prometheus/common/log"
	"math/rand"
	"time"
)

type UserService struct{
}

// 注册
func (s *UserService) Register(mobile, plainpwd, nickname, avatar, sex string) (user model.User, err error) {
	tmp := model.User{}
	_, err = DbEngin.Where("mobile=? ", mobile).Get(&tmp)
	if err != nil {
		return tmp, err
	}

	if tmp.Id > 0 {
		return tmp, errors.New("该账户已经注册")
	}

	tmp.Mobile = mobile
	tmp.Avatar = avatar
	tmp.Nickname = nickname
	tmp.Sex = sex
	tmp.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	tmp.Passwd = util.MakePasswd(plainpwd, tmp.Salt)
	tmp.Createat = time.Now()

	tmp.Avatar = "/static/img/tabbar/me.png"
	// token 随机数
	tmp.Token = fmt.Sprintf("%08d", rand.Int31())
	_, err = DbEngin.InsertOne(&tmp)
	log.Warn()
	return tmp, err
}
