package chat

import (
	"golang.org/x/net/websocket"
)

const SERVER = "golang_chat_bysir"
const SECRET = "y9p0qRr6Yaf~LBBIv3WEJ6c"

type User struct {
	Id      int64
	Name    string
	HeadPic string
	Login   bool
}

type Conn struct {
	ws   *websocket.Conn
	user *User
}


type sendStatus int8

const (
	SS_SUCC      sendStatus = iota + 1 // 发送成功
	SS_NOTONLINE                       // 不在线
)

func (p sendStatus) String() string {
	switch p {
	case SS_SUCC:
		return "SUCCESS"
	case SS_NOTONLINE:
		return "NOTONLINE"
	}
	return "NULL"
}

type Receive struct {
	msgInRow []byte
	conn     *Conn
}
