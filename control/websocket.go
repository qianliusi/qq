package control

// WebChat project main.go
import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	. "qq/model"
	"sync/atomic"
)

var server = NewWebSocketServer()

type WebSocketServer struct {
	UserMap          map[int64]*User
	OnlineUsers      map[*websocket.Conn]*User
	OnlineChan       chan *websocket.Conn
	OfflineChan      chan *websocket.Conn
	BroadCastMsgChan chan *Msg
	nextUserId       int64
}

func NewWebSocketServer() *WebSocketServer {
	s := &WebSocketServer{
		make(map[int64]*User),
		make(map[*websocket.Conn]*User),
		make(chan *websocket.Conn, 10),
		make(chan *websocket.Conn, 10),
		make(chan *Msg, 10),
		0,
	}
	s.Start()
	return s
}
func (s *WebSocketServer) Start() {
	go s.BroadCastMsg()
	go s.UserOnline()
	go s.UserOffline()
}
func (s *WebSocketServer) NextUserId() int64 {
	atomic.AddInt64(&s.nextUserId, 1)
	return s.nextUserId
}
func (s *WebSocketServer) UserOnline() {
	for {
		if c, ok := <-s.OnlineChan; ok {
			log.Printf("new user online: %s->%s\n", c.RemoteAddr().String(), c.LocalAddr().String())
			if _, ok := s.OnlineUsers[c]; !ok {
				u := User{Id: s.NextUserId(), Name: DefaultUserName}
				s.UserMap[u.Id] = &u
				s.OnlineUsers[c] = &u
				msg := Msg{Type: MsgTypeRegister, UserId: u.Id, UserName: u.Name,}
				bytes, _ := json.Marshal(msg)
				err := websocket.Message.Send(c, string(bytes))
				if err != nil {
					log.Printf("send msg error,%s", err.Error())
					s.OfflineChan <- c
				}else {
					s.BroadCastMsgChan <- &Msg{Type: MsgTypeOnline, UserId: u.Id, UserName: u.Name,}
				}
			}
		}
	}
}

func (s *WebSocketServer) UserOffline() {
	for {
		if c, ok := <-s.OfflineChan; ok {
			log.Printf("user offline: %s->%s\n", c.RemoteAddr().String(), c.LocalAddr().String())
			user := s.OnlineUsers[c]
			delete(s.OnlineUsers, c)
			_ = c.Close()
			delete(s.UserMap, user.Id)
			s.BroadCastMsgChan <- &Msg{Type: MsgTypeOffline, UserId: user.Id, UserName: user.Name,}
		}
	}
}

func (s *WebSocketServer) BroadCastMsg() {
	for {
		if m, ok := <-s.BroadCastMsgChan; ok {
			bytes, e := json.Marshal(m)
			if e != nil{
				log.Panicf("json marshal error,%v", m)
				continue
			}
			for k := range s.OnlineUsers {
				err := websocket.Message.Send(k, string(bytes))
				if err != nil {
					log.Printf("send msg error,%s", err.Error())
					s.OfflineChan <- k
				}
			}
		}
	}
}

func (s *WebSocketServer) ModifyUser(user User) {
	if u, ok := s.UserMap[user.Id]; ok {
		u.Name = user.Name
	}
}
func webSocketHandler(conn *websocket.Conn) {
	server.OnlineChan <- conn
	var err error
	for {
		var data string
		err = websocket.Message.Receive(conn, &data)
		if err != nil {
			log.Println("receive err:", err.Error())
			//移除出错的链接
			server.OfflineChan <- conn
			break
		}
		log.Println("receive msg :", data)
		msg := Msg{}
		err = json.Unmarshal([]byte(data), &msg)
		if err != nil {
			log.Println("解析数据异常...", err.Error())
			continue
		}
		msg.UserId=server.OnlineUsers[conn].Id
		if msg.Type == MsgTypeModifyUser {
			server.ModifyUser(User{Id: msg.UserId, Name: msg.UserName})
		}
		server.BroadCastMsgChan <- &msg
	}

}
