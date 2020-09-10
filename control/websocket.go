package control

// WebChat project main.go
import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	. "qq/model"
)

var connUsers = make(map[*websocket.Conn]User)

func webSocketHandler(conn *websocket.Conn) {
	log.Printf("a new ws conn: %s->%s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())
	var err error
	for {
		//判断是否重复连接
		if _, ok := connUsers[conn]; !ok {
			u := User{Name: DefaultUname}
			connUsers[conn] = u
		}
		var data string
		err = websocket.Message.Receive(conn, &data)
		if err != nil {
			fmt.Println("receive err:", err.Error())
			//移除出错的链接
			delete(connUsers, conn)
			break
		}
		fmt.Println("receive msg :", data)

		msg := Msg{}
		err = json.Unmarshal([]byte(data), &msg)
		if err != nil {
			fmt.Println("解析数据异常...", err.Error())
			break
		}
		user := connUsers[conn]
		user.Name = msg.UserName
		for k := range connUsers {
			err := websocket.Message.Send(k, data)
			if err != nil {
				delete(connUsers, k)
			}
		}
	}

}
