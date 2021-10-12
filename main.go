package main

import (
	"fmt"
	"net/http"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

var Sockets []socketio.Conn

func main() {
	var server = socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		AdSockets(s)
		// fmt.Println(Sockets)
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		RemoveSockets(s)
		// fmt.Println(Sockets)
	})

	go server.Serve()
	defer server.Close()

	server.OnEvent("/", "SetName", func(s socketio.Conn, data string) {
		s.SetContext(data)
		AdSockets(s)
		RemoveSockets(s)
		fmt.Println("ctx:", s.Context())
		// fmt.Println(Sockets)
	})

	server.OnEvent("/", "DreamAll", func(s socketio.Conn, data string) {
		// fmt.Println("Dream:", data)
		BroadcastAll("Dream", data)
	})

	server.OnEvent("/", "Dream", func(s socketio.Conn, data string) {
		splited := strings.Split(data, "|")
		context := fmt.Sprintf("%v", s.Context())
		BroadcastOne(splited[0], "Dream", context+"|"+splited[1])
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe("0.0.0.0:8084", nil)
}

func BroadcastAll(event string, data string) {
	for _, socket := range Sockets {
		socket.Emit(event, data)
	}
}

func BroadcastOne(who string, event string, data string) {
	for _, socket := range Sockets {
		if socket.Context() == who {
			socket.Emit(event, data)
		}
	}
}

func remove(slice []socketio.Conn, s int) []socketio.Conn {
	return append(slice[:s], slice[s+1:]...)
}

func indexOf(element socketio.Conn, data []socketio.Conn) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func AdSockets(s socketio.Conn) {
	Sockets = append(Sockets, s)
}

func RemoveSockets(s socketio.Conn) {
	i := indexOf(s, Sockets)
	if i != -1 {
		Sockets = remove(Sockets, i)
	}
}
