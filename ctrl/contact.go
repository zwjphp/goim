package ctrl

import (
	"goim/args"
	"goim/model"
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

func LoadFriend(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	users := contactService.SearchFriend(arg.Userid)
	util.RespOkList(w, users, len(users))
}

func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.Community
	util.Bind(r, &arg)
	if arg.Ownerid == 0 || len(arg.Name) == 0 || len(arg.Icon) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	conn, err := contactService.CreateCommunity(arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, conn, "创建群成功")
	}
}

func JoinCommunity(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 || arg.Dstid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	err := contactService.JoinCommunity(arg.Userid, arg.Dstid)
	// todo 刷新用户的群组信息
	AddGroupId(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "success")
	}
}

func LoadCommunity(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	comunitys := contactService.SearchComunity(arg.Userid)
	util.RespOkList(w, comunitys, len(comunitys))
}


















