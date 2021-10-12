package socket

import (
	"fmt"
	"strings"

	socketio "github.com/googollee/go-socket.io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Sockets []socketio.Conn

func Start(port string) {
	server := socketio.NewServer(nil)

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

	e := echo.New()
	e.HideBanner = true

	e.Static("/", "./public")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Any("/socket.io/", func(context echo.Context) error {
		server.ServeHTTP(context.Response(), context.Request())
		return nil
	})
	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
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
