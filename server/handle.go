package main

import (
	"fmt"
	"net"
	"runtime"
)

func HandleClient(client *client, conn chan net.Conn) {
	go client.Read()
	go client.Write()

	for {
		select {
		case err := <-client.exit:
			fmt.Printf("client has error: %s\n", err.Error())
			client.reConn <- true
			runtime.Goexit()
		case userConn := <-conn:
			user := &userClient{
				conn:  userConn,
				read:  make(chan []byte),
				write: make(chan []byte),
				exit:  make(chan error),
			}
			go user.Read()
			go user.Write()

			go handle(client, user)
		}
	}
}

func handle(client *client, user *userClient) {
	for {
		select {
		case userReceive := <-user.read:
			client.write <- userReceive
		case clientReceive := <-client.read:
			user.write <- clientReceive
		case err := <-user.exit:
			fmt.Println("user has error:", err.Error())
			_ = user.conn.Close()
		case err := <-client.exit:
			fmt.Println("client has error", err.Error())
			_ = client.conn.Close()
			_ = user.conn.Close()
			client.reConn <- true
			runtime.Goexit()
		}
	}
}

func AcceptUserConn(userListener net.Listener, connChan chan net.Conn) {
	userConn, err := userListener.Accept()
	if err != nil {
		panic(err)
	}
	fmt.Printf("user connect:%s\n", userConn.RemoteAddr())
	connChan <- userConn
}
