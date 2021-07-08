package ctrl

import (
	"fmt"
	"goim/model"
	"goim/server"
	"goim/util"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var userService server.UserService

// 登录
func UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		mobile := r.PostForm.Get("mobile")
		plainpwd := r.PostForm.Get("passwd")
		if len(mobile) == 0 || len(plainpwd) == 0 {
			util.RespFail(w, "参数错误")
			return
		}
		user, err := userService.Login(mobile, plainpwd)
		if err != nil {
			util.RespFail(w, err.Error())
		} else {
			tokenString, _ := util.GenToken(user.Mobile)
			util.RespOk(w, user, tokenString)
		}
	}
}

// 注册
func UserRegister(w http.ResponseWriter, r *http.Request) {
	// 1.获取前端传递过来的参数
	if r.Method == http.MethodPost {
		r.ParseForm()
		mobile := r.PostForm.Get("mobile")
		plainpwd := r.PostForm.Get("passwd")
		uuid := r.PostForm.Get("uuid")
		code := r.PostForm.Get("code")
		log.Print(mobile)
		log.Print(plainpwd)
		log.Print(uuid)
		log.Print(code)
		if len(mobile) == 0 || len(plainpwd) == 0 || len(uuid) == 0 {
			util.RespFail(w, "参数错误")
			return
		}

		err := util.CaptchaVerifyHandle(uuid, code)
		if err != nil {
			util.RespFail(w, err.Error())
			return
		}
		rand.Seed(time.Now().UnixNano()) // 设置种子数为当前时间
		nickname := fmt.Sprintf("user%06d", rand.Int31())
		avatar := ""
		sex := model.SEX_UNKNOW
		user, err := userService.Register(mobile, plainpwd, nickname, avatar, sex)
		if err != nil {
			util.RespFail(w, err.Error())
		} else {
			util.RespOk(w, user, "")
		}
	}
}

func GetCaptcha(w http.ResponseWriter, r *http.Request) {
	util.GenerateCaptchaHandler(w, r)
}

func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := r.PostForm.Get("userid")
	avatar := r.PostForm.Get("avatar")
	if len(userid) == 0 || len(avatar) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	id, _ := strconv.Atoi(userid)
	userService.UserInfo(int64(id), avatar)
	util.RespOk(w, nil, "")
}





















