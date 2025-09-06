package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var clients []net.Conn

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		clients = append(clients, conn)
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		broadcast(msg)
	}
}

func broadcast(msg string) {
	fmt.Println(msg)

	for _, client := range clients {
		client.Write([]byte(msg))
	}
}
