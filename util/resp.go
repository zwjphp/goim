package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type H struct {
	Code int          `json:"code"`
	Msg string        `json:"msg"`
	Data interface{}  `json:"data,omitempty"`
	Rows interface{}  `json:"rows,omitempty"`
	Total interface{} `json:"total,omitempty"`
}

// 错误返回信息体
func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

func RespOk(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0,  data,  msg)
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	// 将结构体转化成json字符串
	ret, err := json.Marshal(h)
	if err != nil {
		log.Printf(err.Error())
	}
	w.Write(ret)
}

// 返回列表信息 及 总条数
func RespOkList(w http.ResponseWriter, lists interface{}, total interface{}) {
	RespList(w, 0, lists, total)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}

	ret, err := json.Marshal(h)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(ret)
}


















