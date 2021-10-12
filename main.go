package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	socketio "github.com/googollee/go-socket.io"

	"github.com/gin-gonic/gin"
)

var Sockets []socketio.Conn

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
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

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("./home", gin.Dir("./public", false))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	if err := router.Run("0.0.0.0:8084"); err != nil {
		fmt.Println("failed run app: ", err)
	}
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
