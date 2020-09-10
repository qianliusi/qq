package control

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	. "qq/config"
	"qq/service"
	"qq/utils"
	"runtime/debug"
	"strconv"
)

func Run() {
	mux := http.NewServeMux()

	mux.Handle("/", http.RedirectHandler("/static/room.html", 302))

	mux.HandleFunc("/static/", safeHandler(staticHandler))
	mux.Handle("/ws", websocket.Handler(webSocketHandler))
	mux.HandleFunc("/getAllRooms", safeHandler(getAllRooms))
	mux.HandleFunc("/getRoomInfo", safeHandler(getRoomInfo))

	//http.HandleFunc("/static/", safeHandler(staticHandler))
	//http.Handle("/ws", websocket.Handler(webSocketHandler))
	//http.HandleFunc("/getAllRooms", safeHandler(getAllRooms))
	//http.HandleFunc("/getRoomInfo", safeHandler(getRoomInfo))
	err := http.ListenAndServe("", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}

}

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				log.Printf("WARN: panic in %v - %v\n", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, r.URL.Path)
}

func getAllRooms(w http.ResponseWriter, r *http.Request) {
	utils.SendResponseJson(w, service.GetAllRooms())
}

func getRoomInfo(w http.ResponseWriter, r *http.Request) {
	roomId, _ := strconv.Atoi(r.FormValue("roomId"))
	utils.SendResponseJson(w, service.GetRoomInfo(roomId))
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	file := "." + r.URL.Path
	if exists := utils.Exists(file); !exists {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, file)
}

func staticDirHandler(mux *http.ServeMux, prefix string, flags int) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := StaticDir + r.URL.Path[len(prefix)-1:]
		if (flags & ListDir) == 0 {
			if exists := utils.Exists(file); !exists {
				http.NotFound(w, r)
				return
			}
		}
		http.ServeFile(w, r, file)
	})
}
