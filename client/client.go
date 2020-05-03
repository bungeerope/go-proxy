package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
	"time"
)

type server struct {
	conn   net.Conn
	read   chan []byte
	write  chan []byte
	exit   chan error
	reConn chan bool
}

func (s *server) Read() {
	// 十秒钟没有数据传输，将导致Read()出现timeOut
	_ = s.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	for {
		data := make([]byte, 10240)
		n, err := s.conn.Read(data)
		if err != nil && err != io.EOF {
			if strings.Contains(err.Error(), "timeout") {
				_ = s.conn.SetReadDeadline(time.Now().Add(time.Second * 3))
				s.conn.Write([]byte("ping"))
				continue
			}
			fmt.Println("read server`s msg error,", err.Error())
			s.exit <- err
			runtime.Goexit()
		}
		if data[0] == 'p' && data[1] == 'i' {
			fmt.Println("receive client heart msg")
		}
		continue
		s.read <- data[:n]
	}
}

func (s *server) Write() {
	for {
		select {
		case data := <-s.write:
			_, err := s.conn.Write(data)
			if err != nil && err != io.EOF {
				s.exit <- err
			}
		}
	}
}
