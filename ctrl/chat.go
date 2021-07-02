package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"goim/model"
	"goim/server"
	"goim/validates"
	"net/http"
	"strconv"
	"gopkg.in/fatih/set.v0"
	"sync"
)

// 本核心在与形成userId 和 Node的映射关系
type Node struct {
	Conn *websocket.Conn
	DataQueue chan []byte // 并行转串行
	GroupSets set.Interface
}

// 映射关系
var clientMap = make(map[int64]*Node, 0)

// 读写锁
var rwlocker sync.RWMutex
var log = logrus.New()

func Chat(w http.ResponseWriter, r *http.Request) {
	// TODO校验token合法性
	// checkToken
	// 获取路由参数
	query     := r.URL.Query()
	id        := query.Get("id")
	token     := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64) // 将字符串转换为int64类型
	isvalida  := checkToken(userId, token)

	// 如果isvalida=true
	// isvalida=false
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(w, r, nil)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"animal": "walrus",
			"size": 10,
		}).Warn(err.Error())
		return
	}

	// TODO 获得 conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// TODO 获取用户全部群id
	comIds := contactService.SearchComunityIds(userId)
	for _, v := range comIds {
		node.GroupSets.Add(v)
	}

	// todo userid 和 node 形成绑定关系
	rwlocker.Lock() // 锁
	clientMap[userId] = node
	rwlocker.Unlock() // 释放锁
	// todo 完成发送逻辑
	go sendproc(node)
	// todo 完成接收逻辑
	go recvproc(node)
}

// 发送协程(写)
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.WithFields(logrus.Fields{
					"animal": "walrus",
					"size": 10,
				}).Warn(err.Error())
				return
		    }
		}
	}
}

// 接收协程(读)
func recvproc(node *Node) {
	defer node.Conn.Close()
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.WithFields(logrus.Fields{
				"animal": "walrus",
				"size": 10,
			}).Warn(err.Error())
			return
		}
		// todo 对data进一步处理
		dispatch(data)
		fmt.Printf("recv <= %s\n", data)
	}
}

// 调度逻辑处理
func dispatch(data []byte) {
	// TODO 解析data为message
	msg := model.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Warn(err.Error())
		return
	}
	b := validates.VerificationFilter(msg.Content)
	if b {
		msg.Cmd = model.CMD_FILTER
	}
	// TODO 根据cmd对逻辑进行处理
	switch msg.Cmd {
	case model.CMD_HEART:
		// 心跳
	case model.CMD_SINGLE_MSG:
		sendMsg(msg.Dstid, data)
		server.Rpush("chat_10", data)
	case model.CMD_ROOM_MSG:
		// 群聊
		for _, v := range clientMap {
			if v.GroupSets.Has(msg.Dstid) {
				v.DataQueue <- data
			}
		}
		server.Rpush("chat_11", data)
	case model.CMD_QUIT:
		// 退出
		DelClientMapID(msg.Userid)
	case model.CMD_NEW_FRIEND:
		// 通知新朋友添加
		sendMsg(msg.Dstid, data)
	case model.CMD_FILTER:
		sendMsg(msg.Userid, []byte(`{"dstid":`+strconv.FormatInt(msg.Dstid,10)+`,"cmd":`+strconv.Itoa(model.CMD_FILTER)+`}`))
	}
}

// 用户退出删除连接
func DelClientMapID(userId int64) {
	rwlocker.Lock()
	_, ok := clientMap[userId]
	if ok {
		delete(clientMap, userId)
	}
	rwlocker.Unlock()
}

// 发送消息
func sendMsg(userId int64, msg []byte) {
	rwlocker.RLock() // 锁
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

func checkToken(userId int64, token string) bool {
	// 从数据库中查询 并 比对
	user := userService.Find(userId)
	return user.Token == token
}