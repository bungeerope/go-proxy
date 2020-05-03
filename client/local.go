package main

import (
	"io"
	"net"
)

type local struct {
	conn  net.Conn
	read  chan []byte
	write chan []byte
	exit  chan error
}

func (l *local) Read() {
	for {
		data := make([]byte, 10240)
		n, err := l.conn.Read(data)
		if err != nil && err != io.EOF {
			l.exit <- err
		}
		l.read <- data[:n]
	}
}

func (l *local) Write() {
	for {
		select {
		case data := <-l.write:
			_, err := l.conn.Write(data)
			if err != nil && err != io.EOF {
				l.exit <- err
			}
		}
	}
}
