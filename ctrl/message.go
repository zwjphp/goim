package ctrl

import (
	"encoding/json"
	"goim/args"
	"goim/model"
	"goim/server"
	"goim/util"
	"net/http"
	"time"
)

var messageService server.MessageService

func ChatHistory(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 || arg.Dstid == 0 || arg.Cmd == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	var chat string
	if arg.Cmd == model.CMD_ROOM_MSG {
		chat = "chat_11"
	} else {
		chat = "chat_10"
	}
	LRange, b2 := server.LRange(chat, int64(arg.GetPageSize()), int64(arg.GetPageSize()))
	if b2 {
		var msg model.Message
		msgList := make([]model.Message, 0)
		for _, data := range LRange{
			json.Unmarshal([]byte(data), &msg)
			if arg.Cmd == model.CMD_ROOM_MSG {
				if arg.Userid != 0 && arg.Dstid == msg.Dstid {
					msg.Createat = time.Now().Unix()
					msgList = append([]model.Message{msg}, msgList...)
				}
			} else {
				if msg.Userid != 0 && (msg.Dstid == arg.Userid && msg.Userid == arg.Dstid) || (arg.Userid == msg.Userid && arg.Dstid == msg.Dstid) {
					msg.Createat = time.Now().Unix()
					msgList = append([]model.Message{msg}, msgList...)
				}
			}
		}
		util.RespOkList(w, msgList, len(msgList))
		return
	} else {
		history := messageService.GetChatHistory(arg.Userid, arg.Dstid, arg.Cmd, arg.GetPageSize(), arg.GetPageSize())
		util.RespOkList(w, history, len(history))
	}
}
