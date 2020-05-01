package main

import (
	"fmt"
	"net"
)

func handle(server *server) {
	data := <-server.read
	localConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		panic(err)
	}
	local := &local{
		conn:  localConn,
		read:  make(chan []byte),
		write: make(chan []byte),
		exit:  make(chan error),
	}
	go local.Read()
	go local.Write()

	local.write <- data

	for {
		select {
		case data := <-server.read:
			local.write <- data
		case data := <-local.read:
			server.write <- data
		case err := <-server.exit:
			fmt.Printf("server have a error: %s", err.Error())
			_ = server.conn.Close()
			_ = local.conn.Close()
			server.reConn <- true
		case err := <-local.exit:
			fmt.Printf("local have a error: %s", err.Error())
			_ = local.conn.Close()
		}
	}

}
