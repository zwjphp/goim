package ctrl

import (
	"goim/args"
	"goim/server"
	"goim/util"
	"goim/validates"
	"net/http"
)

var contactService server.ContactService
var contactValidate validates.ContactValidate

func Addfriend(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	contactValidates, err := contactValidate.ContactValidates(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, contactValidates)
		return
	}
	// 调用service
	err = contactService.AddFriend(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "好友添加成功")
	}
}