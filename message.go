package chat

import "encoding/json"

type msgType int
type msgTo int

const (
	MTP_TEXT        msgType = iota + 1
	MTP_IMAGE
	MTP_AUDIO
	MTP_LOGIN        // login

	// out

	MTP_NY_MSG_FAIL  // msg not send
	MTP_NY_MSG_READ  // msg read
	MTP_NY_MSG_SUCC  // msg send success
	MTP_LOGIN_FAIL
	MTP_LOGIN_SUCC
)

const (
	MTO_RSER msgTo = iota + 1
	MTO_ROOM
)

// server receive message from user
type MessageIn struct {
	Types msgType `json:"type,emitempty"` // 1 text, 2 image, 3 audio
	Data  string `json:"data,emitempty"`  // text or file path
	To    msgTo `json:"to,emitempty"`     // 1 user,2 room
	ToId  int64 `json:"to_id,emitempty"`  // user id or room id
	Id    int64 `json:"id,emitempty"`     // id应该尽量不一样, 用于通知发送了的msg的状态
}

func (p *MessageIn) Encode() []byte {
	d, _ := json.Marshal(p)
	return d
}

func (p *MessageIn) Decode(d []byte) (error) {
	err := json.Unmarshal(d, p)
	return err
}

// -------- out ----------

// server send message to user
type MessageOut struct {
	Types     msgType `json:"type,emitempty"`     // 1 text, 2 image, 3 audio
	FromId    int64 `json:"from,emitempty"`       // from uid
	FromPic   string `json:"form_pic,emitempty"`  // from headPic
	FromName  string `json:"form_name,emitempty"` // from name
	Timestamp int64 `json:"timestamp,emitempty"`  // send time
	Data      string `json:"data,emitempty"`      // text or file path
	Id        int64 `json:"id,emitempty"`         // 这里的id是发送者发送的id
}

func (p *MessageOut) Encode() []byte {
	d, _ := json.Marshal(p)
	return d
}

func (p *MessageOut) Decode(d []byte) (error) {
	err := json.Unmarshal(d, p)
	return err
}
