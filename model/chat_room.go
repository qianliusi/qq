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

const DefaultUname = "匿名"

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ChatRoomId int    `json:"chatRoomId"`
}

type Msg struct {
	Id       int    `json:"id"`
	Type     string `json:"type"`
	UserName string `json:"userName"`
	Content  string `json:"content"`
}

type Sender interface {
	Send() error
}
