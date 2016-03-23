package chat

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
    MAXWAIT = 10
	SETNAMEFLAG = "setname:"
	TIMEOUT = 20
)

type Server struct {
	running  bool    // running flag
	listener net.Listener
	client   map[net.Conn]*Client   //connected clients
	broadcast MESSAGE    //recieved msg from client
    join    chan net.Conn
    quit    chan net.Conn
	existname map[string]struct{}  //client namespace
	weather map[string]string //simple weather cache
}

// New a server, then listen event
func NewServer(addr string) (*Server, error) {
	s := &Server{
	    client: make(map[net.Conn]*Client),
		broadcast: make(MESSAGE),
		join: make(chan net.Conn, MAXWAIT),
		quit: make(chan net.Conn),
		existname: make(map[string]struct{}),
		weather: make(map[string]string),
	}
	netProto := "tcp"
	if strings.Contains(netProto, "/") {
		netProto = "uinx"
	}
	var err error
	//listen local network
	s.listener, err = net.Listen(netProto, addr)
	if err != nil {
		return nil, err
	}
	s.ListenEvent()
	s.WatchLoop()
	return s, nil
}

// accept a connect
func (s *Server) Run() {
	s.running = true
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			//log error
			fmt.Printf("accept error, reason:%v\n", err.Error())
			continue
		}
		fmt.Printf("client %v connected\n", conn.RemoteAddr().String())
		s.join <- conn
	}
}

// watch every client
func (s *Server) WatchLoop() {
	go func() {
		for {
			for conn, client := range s.client {
				fmt.Printf("client last time:%v\n", client.lastact)
				if int(time.Now().Sub(client.lastact).Seconds()) >= TIMEOUT {
					s.removeFromClient(conn)
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

// listen kinds of events
func (s *Server) ListenEvent() {
	go func() {
		for {
			select {
				//broadcast
			//case message := <-s.broadcast:
			//	s.broacastMessage(message)
				//join
			case conn := <-s.join:
				fmt.Printf("rcv a connect %v\n", conn.RemoteAddr().String())
				s.join2Client(conn)
				//remove
			case conn := <-s.quit:
			    s.removeFromClient(conn)
			}
		}
	}()
}

func (s *Server) ClientInit(client Client) {
}

// join a Clien into Server.client map
func (s *Server) join2Client(conn net.Conn) {
    client := NewClient(conn)
	s.client[conn] = client
	var first = true
	// first send set name response
	client.outgoing<- "Set name first!format:setname:xxx"
	fmt.Printf("new a client name:%v\n", client.name)

	// quit event
    go func() {
	    for {
		    quit := <-client.quit
			_, found := s.existname[client.name]
			if !found {
				continue
			}
			fmt.Printf("client quit,name:%v\n", client.name)
			//send to server
			s.quit<-quit
		}
	}()

	// client send message
    go func() {
	    for {
		    msg := <-client.incoming
			client.lastact = time.Now()
			if first {
				if strings.HasPrefix(msg, SETNAMEFLAG) {
					clientname := strings.Split(msg, ":")[1]
					_, found := s.existname[clientname]
					if found {
						client.outgoing<- fmt.Sprintf("name %v has be used.", clientname)
						continue
					}
					client.SetName(clientname)
					s.existname[clientname] = struct{}{}
					client.outgoing<- fmt.Sprintf("Set name success %v!", clientname)
					client.outgoing <- fmt.Sprintf("Welcome!people online:%v", len(s.client))
					s.broacastMessage(fmt.Sprintf("%v 已上线", clientname), clientname)
					first = false
				} else if msg == "\\quit"{
					s.removeFromClient(conn)
					//client.quit <- conn
					//continue
				} else {
					client.outgoing <- fmt.Sprintf("Set name first!format:setname:xxx")
				}
				//扩展其他命令,such as weather,ticket......
			} else if strings.HasPrefix(msg, "\\"){
				if msg == "\\quit" {
					//s.removeFromClient(conn)
					client.quit <- conn
				} else if strings.HasPrefix(msg, "\\天气") {
				    city := strings.Split(msg, " ")[1]
					value, found := s.weather[city]
					if found {
					    client.outgoing <- value
						continue
					}
					weather, err := s.GetWeather(city)
					if err != nil {
						// 不考虑cache变了之后的更新，需要加cache watch功能
						// 2小时之后自动更新cache
						s.weather[city] = weather
						client.outgoing <- fmt.Sprintf("get weather failed, reason:%v", err.Error())
					} else {
					    client.outgoing <- weather
					}
				}
			}else {
				saywhat := fmt.Sprintf("%v;%v", client.name, msg)
				fmt.Printf("%v\n", saywhat)
				//send to sever
				//s.broadcast <- saywhat
				s.broacastMessage(saywhat, client.name)
			}
		}
	}()
	fmt.Printf("join done\n")
}

// remove a client from client map
func (s *Server) removeFromClient(conn net.Conn) {
	var quitname = s.client[conn].name
	if conn != nil {
		fmt.Printf("remove from client list %v\n", conn.RemoteAddr().String())
		s.client[conn].Close()
		delete(s.client, conn)
		delete(s.existname, quitname)
	}
	s.broacastMessage(fmt.Sprintf("%v 已下线", quitname), "")
}


// close server
func (s *Server) Close() error {
	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}
	return nil
}

// broadcast msg to every connected client
func (s *Server) broacastMessage(message string, except string) {
    for _, client := range s.client {
		if client.name != except {
			client.outgoing <- message
		}
	}
}

