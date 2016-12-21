package chat

import (
	"fmt"
	"github.com/bysir-zl/bjson"
	"github.com/bysir-zl/bygo/util"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

type Server struct {
	hub hub
}

func (p *Server) echoHandler(ws *websocket.Conn) {
	// read one message use to login
	var firstMsgChan chan []byte = make(chan []byte)

	go func() {
		msg := make([]byte, 1024 * 256)
		n, err := ws.Read(msg)
		if err != nil {
			firstMsgChan <- nil
			return
		}
		firstMsgChan <- msg[:n]
	}()

	select {
	case <-time.After(5 * time.Second):
		close(firstMsgChan)
		ws.Write(util.S2B(fmt.Sprintf(`{"type":%d,"data":"timeout"}`, MTP_LOGIN_FAIL)))
		return
	case firstMsg := <-firstMsgChan:
		if firstMsg == nil {
			ws.Write(util.S2B(fmt.Sprintf(`{"type":%d,"data":"can not read"}`, MTP_LOGIN_FAIL)))
			return
		}
		msgIn := MessageIn{}
		err := msgIn.Decode(firstMsg)
		if err != nil {
			ws.Write(util.S2B(fmt.Sprintf(`{"type":%d,"data":"format error"}`, MTP_LOGIN_FAIL)))
			return
		}
		if msgIn.Types != MTP_LOGIN {
			ws.Write(util.S2B(fmt.Sprintf(`{"type":%d,"data":"not log action"}`, MTP_LOGIN_FAIL)))
			return
		}
	// decode data
		bj, err := bjson.New(util.S2B(msgIn.Data))
		if err != nil {
			ws.Write(util.S2B(fmt.Sprintf(`{"type":%d,"data":"format error"}`, MTP_LOGIN_FAIL)))
			return
		}
		token := bj.Pos("token").String()
		uid, err := VerifyUser(token)
		if err != nil {
			ws.Write(util.S2B(fmt.Sprintf(`{"type":%d,"data":"auth error, %v"}`, MTP_LOGIN_FAIL,err)))
			return
		}

		ws.Write(util.S2B(fmt.Sprintf(`{"type":%d,"data":"%d"}`, MTP_LOGIN_SUCC, uid)))

		headPic := bj.Pos("head_pic").String()
		name := bj.Pos("name").String()
		user := &User{
			Id:     uid,
			Name:   name,
			HeadPic:headPic,
		}

		c := &Conn{
			ws:  ws,
			user:user,
		}
		p.hub.onConnect <- c
		defer func(c *Conn) {
			p.hub.onClose <- c
		}(c)

		msg := make([]byte, 1024 * 256)
		for {
			// for receive
			n, err := ws.Read(msg)
			if err != nil {
				return
			}
			p.hub.onReceive <- &Receive{
				conn:    c,
				msgInRow:msg[:n],
			}
		}

	}
}

func (p *Server) Server() error {
	p.hub = newHub()
	p.hub.Start()
	http.Handle("/chat", websocket.Handler(p.echoHandler))
	return http.ListenAndServe(":11000", nil)
}
