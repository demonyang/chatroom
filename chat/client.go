package chat

import (
    "bufio"
	"net"
	"fmt"
	"time"
)

type MESSAGE chan string

type Client struct {
    name string
	conn net.Conn
	incoming  MESSAGE
	outgoing  MESSAGE
	quit  chan net.Conn
	read   *bufio.Reader
	write  *bufio.Writer
	latesttime time.Time
}

func NewClient(con net.Conn) *Client {
    client := &Client{
	    conn: con,
		name: con.RemoteAddr().String(),
		incoming: make(MESSAGE),
		outgoing: make(MESSAGE),
		quit: make(chan net.Conn),
		read: bufio.NewReader(con),
		write: bufio.NewWriter(con),
		latesttime: time.Now(),
	}
	//listen
	fmt.Printf("start listen\n")
	client.Listen()
	fmt.Printf("client listen done\n")
	return client
}

func (c *Client) Listen() {
	go func() {
		for{
			if line, _, err := c.read.ReadLine(); err == nil {
				c.incoming <- string(line)
			} else {
				fmt.Printf("read err:%v\n", err.Error())
				c.quitc()
				return
			}
		}
	}()

	go func() {
	    for data := range c.outgoing {
		    if _, err := c.write.WriteString(data+"\n"); err != nil {
			    c.quitc()
				return
			}
			if err := c.write.Flush(); err != nil {
			    fmt.Printf("write err %v", err.Error())
				c.quitc()
				return
			}
		}
	}()
}

func (c *Client) Putout(str string) {
    c.outgoing<-str
}

func (c *Client) Getin() string {
    return <-c.incoming
}

func (c *Client) quitc() {
	c.quit <- c.conn
}

func (c *Client) Close() {
    if c.conn != nil {
	    c.conn.Close()
	}
}

func (c *Client) SetName(newname string) {
    c.name = newname
}
