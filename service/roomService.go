package service

import (
	. "qq/config"
	. "qq/dao"
	. "qq/model"
)

var rooms Rooms

func init() {
	initAllRooms()
}

func GetAllRooms() Rooms {
	return rooms
}

func initAllRooms() {
	CommonDao.Read(RoomsPath, &rooms)
}

func GetRoomInfo(roomId int) *ChatRoom {
	for e := range rooms.Rooms {
		if rooms.Rooms[e].Id == roomId {
			return &rooms.Rooms[e]
		}
	}
	return nil
}
