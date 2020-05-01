package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	host       string
	localPort  int
	remotePort int
)

func init() {
	flag.StringVar(&host, "h", "127.0.0.1", "remote server ip")
	flag.IntVar(&localPort, "l", 8080, "local port")
	flag.IntVar(&remotePort, "r", 3333, "remote port")
}

func main() {
	flag.Parse()
	target := net.JoinHostPort(host, fmt.Sprintf("%d", remotePort))
	for {
		serverConn, err := net.Dial("tcp", target)
		if err != nil {
			panic(err)
		}
		fmt.Printf("connected server:%s \n", serverConn.RemoteAddr())
		server := &server{
			conn:   serverConn,
			read:   make(chan []byte),
			write:  make(chan []byte),
			exit:   make(chan error),
			reConn: make(chan bool),
		}
		go server.Read()
		go server.Write()

		go handle(server)
		<-server.reConn
		_ = server.conn.Close()
	}
}
