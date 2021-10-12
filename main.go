package main

import (
	s "gobackmessage/socket"
)

func main() {
	port := "8084"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	s.Start(port)
}
