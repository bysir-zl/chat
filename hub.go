package chat

import (
	"github.com/bysir-zl/bygo/log"
	"sync"
	"time"
	"github.com/bysir-zl/bygo/util"
)

type hub struct {
	connects map[int64]*Conn // userid : conn

	onReceive chan *Receive
	onClose   chan *Conn
	onConnect chan *Conn
	stopChan  chan int
}

func newHub() hub {
	return hub{
		connects:    map[int64]*Conn{},
		onReceive:   make(chan *Receive, 10240),
		onClose:     make(chan *Conn, 1024),
		onConnect:   make(chan *Conn, 1024),
	}
}

var connsLook sync.RWMutex

func (p *hub) Start() {
	go func() {
		for {
			select {
			case receive := <-p.onReceive:
				go p.handleReceive(receive)
			case conn := <-p.onConnect:
				connsLook.Lock()
				p.connects[conn.user.Id] = conn
				connsLook.Unlock()
			case conn := <-p.onClose:
				delete(p.connects, conn.user.Id)
			}
		}
		p.stopChan <- 0
	}()
}

func (p *hub) handleReceive(r *Receive) {
	msgIn := new(MessageIn)
	err := msgIn.Decode(r.msgInRow)
	if err != nil {
		log.Error("hub handlemessage", err, util.B2S(r.msgInRow))
		return
	}
	user := r.conn.user
	switch msgIn.To {
	case MTO_RSER:
		msgOut := &MessageOut{
			Types:    msgIn.Types,
			Data:     msgIn.Data,
			Id:       msgIn.Id,
			FromId:   user.Id,
			FromName: user.Name,
			FromPic:  user.HeadPic,
			Timestamp:time.Now().Unix(),
		}
		err, status := p.sendToUser(msgIn.ToId, msgOut)
		// notify
		var types msgType
		var data string
		if err != nil {
			types = MTP_NY_MSG_FAIL
			data = err.Error()
		} else {
			types = MTP_NY_MSG_SUCC
			data = status.String()
		}
		e, s := p.notifyUser(user.Id, types, data, msgIn.Id)
		if e != nil {
			log.Error("send notify", e)
		} else if s != 1 {
			log.Error("send notify", "notify user failed ", s)
		}
	case MTO_ROOM:
	// todo room
	}

}

// status 1:success, 2:not online
func (p *hub) sendToUser(toUid int64, msg *MessageOut) (err error, status sendStatus) {
	connsLook.RLock()
	conn, ok := p.connects[toUid]
	connsLook.RUnlock()
	if !ok {
		status = SS_NOTONLINE
		// todo save to db
		// or send a save db job to worker
		return
	}
	_, err = conn.ws.Write(msg.Encode())
	status = SS_SUCC
	return
}

func (p *hub) notifyUser(uid int64, types msgType, data string, msgId int64) (err error, status sendStatus) {
	msgOut := &MessageOut{
		Types:    types,
		Data:     data,
		Id:       msgId,
		FromId:   1,
		FromName: "SYS",
		FromPic:  "",
		Timestamp:time.Now().Unix(),
	}
	return p.sendToUser(uid, msgOut)
}
