package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type client struct {
	conn   net.Conn
	read   chan []byte
	write  chan []byte
	exit   chan error
	reConn chan bool
}

func (c *client) Read() {
	_ = c.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	for {
		data := make([]byte, 10240)
		n, err := c.conn.Read(data)
		if err != nil && err != io.EOF {
			if strings.Contains(err.Error(), "timeout") {
				_ = c.conn.SetReadDeadline(time.Now().Add(time.Second * 3))
				c.conn.Write([]byte("pi"))
				continue
			}
			fmt.Println("read msg error")
			c.exit <- err
		}
		if data[0] == 'p' && data[1] == 'i' {
			fmt.Println("server receive heart msg")
			continue
		}
		c.read <- data[:n]
	}
}

func (c *client) Write() {
	for {
		select {
		case data := <-c.write:
			_, err := c.conn.Write(data)
			if err != nil && err != io.EOF {
				c.exit <- err
			}
		}
	}
}
