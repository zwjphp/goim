package handler

import (
	"github.com/gorilla/mux"
	"goim/ctrl"
	"goim/middleware"
	"goim/util"
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
	indexRouter.HandleFunc("/auth", util.AuthHandler).Methods(http.MethodPost)  // 获取token

	// 文件
	fileRouter := r.PathPrefix("/attach").Subrouter()
	fileRouter.HandleFunc("/upload", ctrl.UploadLocal).Methods(http.MethodPost, http.MethodOptions) // 文件上传

	indexRouter.HandleFunc("/chat", ctrl.Chat) // ws

	// 需要验证Token
	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.JWTAuthMiddleware, middleware.AccessLogging)

	authRouter.HandleFunc("/contact/addfriend", ctrl.Addfriend).Methods(http.MethodPost) // 添加好友

}
