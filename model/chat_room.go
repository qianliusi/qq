package model

type ChatRoom struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	users []User `json:"users"`
}

type Rooms struct {
	Rooms []ChatRoom `json:"rooms"`
}

type Users struct {
	Users []User `json:"users"`
}

const DefaultUserName = "游客"

type User struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	ChatRoomId int    `json:"chatRoomId"`
}

const (
	MsgTypeRegister   = "register"
	MsgTypeSend       = "send"
	MsgTypeOnline     = "online"
	MsgTypeOffline    = "offline"
	MsgTypeModifyUser = "modifyUser"
)

type Msg struct {
	Id       int    `json:"id"`
	Type     string    `json:"type"`
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	Content  string `json:"content"`
}

type Sender interface {
	Send() error
}
