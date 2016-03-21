package main

import (
	"fmt"
	"flag"
    "demon/chatroom/chat"
)

var address = flag.String("addr", "", "Server listend ip and port")

func main() {
	flag.Parse()
	if *address == "" {
	    fmt.Println("Address can't be null")
		return
	}
	s, err := chat.NewServer(*address)
	if err != nil {
		fmt.Errorf("new server failed:%v", err.Error())
	}
	s.Run()
}
