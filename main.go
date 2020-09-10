package main

import (
	"fmt"
	"qq/control"
	"qq/dao"
	"qq/service"
)

func main() {
	control.Run()
	//test()
}

func test() {
	dao.Test()
	fmt.Println(service.GetAllRooms())

}
