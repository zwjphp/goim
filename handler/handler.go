package handler

import (
	"github.com/gorilla/mux"
	"goim/ctrl"
	"net/http"
)

func RegisterRoutes(r *mux.Router)  {
	// 1. 提供静态资源目录支持
	r.Handle("/apidoc/", http.FileServer(http.Dir(".")))
	r.Handle("/mnt/", http.FileServer(http.Dir(".")))

	// 不需要验证
	indexRouter := r.PathPrefix("/index").Subrouter()
	// 绑定请求的处理函数
	indexRouter.HandleFunc("/getCaptcha", ctrl.GetCaptcha).Methods(http.MethodPost)  // 获取验证码
	indexRouter.HandleFunc("/user/register", ctrl.UserRegister).Methods(http.MethodPost) // 注册
	indexRouter.HandleFunc("/user/login", ctrl.UserLogin).Methods(http.MethodPost) // 登录

}
