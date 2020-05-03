package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	localPort  int
	remotePort int
)

func init() {
	flag.IntVar(&localPort, "p", 8081, "local port")
	flag.IntVar(&remotePort, "r", 3333, "remote port")
}

func main() {
	flag.Parse()
	// 处理异常信息
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	clientListener, err := net.Listen("tcp", fmt.Sprintf(":%d", remotePort))
	if err != nil {
		panic(err)
	}
	fmt.Printf("listen remote port %d ,waitting for client connectting\n", remotePort)
	userListener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		panic(err)
	}
	fmt.Printf("listen local port %d ,waitting for user connectting\n", localPort)
	for {
		clientConn, err := clientListener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Printf("client %s is connected", clientConn.RemoteAddr())

		client := &client{
			conn:   clientConn,
			read:   make(chan []byte),
			write:  make(chan []byte),
			exit:   make(chan error),
			reConn: make(chan bool),
		}
		userConnChan := make(chan net.Conn)
		go AcceptUserConn(userListener, userConnChan)
		go HandleClient(client, userConnChan)
		<-client.reConn
		fmt.Println("It`s reconnecting ...")
	}
}
