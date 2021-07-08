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

// 登录
func (s *UserService) Login(mobile, plainpwd string) (user model.User, err error) {
	tmp := model.User{}
	_, err = DbEngin.Where("mobile=?", mobile).Get(&tmp)
	if err != nil {
		return tmp, err
	}
	// 账号不存在
	if tmp.Id == 0 {
		return tmp, errors.New("账号不存在")
	}
	// 检测密码
	if !util.ValidatePasswd(plainpwd, tmp.Salt, tmp.Passwd) {
		return tmp, errors.New("密码不正确")
	}
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := util.Md5Encode(str)
	tmp.Token = token
	DbEngin.ID(tmp.Id).Cols("token").Update(&tmp)
	return tmp, nil
}

// 查询某个用户信息
func (s *UserService) Find(userId int64) (user model.User) {
	tmp := model.User{}
	DbEngin.ID(userId).Get(&tmp)
	return tmp
}

// 更新用户数据
func (s *UserService) UserInfo(userId int64, avatar string) {
	tmp := model.User{}
	tmp.Avatar = avatar
	DbEngin.ID(userId).Cols("avatar").Update(&tmp)
}
















