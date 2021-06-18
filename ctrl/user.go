package ctrl

import (
	"fmt"
	"goim/model"
	"goim/server"
	"goim/util"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var userService server.UserService
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
