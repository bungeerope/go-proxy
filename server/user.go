package main

import (
	"io"
	"net"
	"time"
)

type userClient struct {
	conn  net.Conn
	read  chan []byte
	write chan []byte
	exit  chan error
}

func (u *userClient) Read() {
	_ = u.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	for {
		data := make([]byte, 10240)
		n, err := u.conn.Read(data)
		if err != nil && err != io.EOF {
			u.exit <- err
		}
		u.read <- data[:n]
	}
}

func (u *userClient) Write() {
	for {
		select {
		case data := <-u.write:
			_, err := u.conn.Write(data)
			if err != nil && err != io.EOF {
				u.exit <- err
			}
		}
	}
}
