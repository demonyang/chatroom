package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"flag"
	"demon/chatroom/chat"
)

var svr = flag.String("svr", "", "The  server to connect")

func main() {
	flag.Parse()
	if *svr == "" {
	    fmt.Println("server address can't be null")
		return
	}
	conn, err := net.Dial("tcp", *svr)
	if err != nil {
		fmt.Errorf("connect to %v error, raason:%v", *svr, err.Error())
		return
	}
	defer conn.Close()
	client := chat.NewClient(conn)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	go func() {
		for {
			out.WriteString(client.Getin() + "\n")
			out.Flush()
			//fmt.Printf("write msg from write\n")
		}
	}()
	for {
		line, _, _ := in.ReadLine()
		//fmt.Printf("read msg from read\n")
		client.Putout(string(line))
	}
}
