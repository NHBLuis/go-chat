package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Client struct {
	connection net.Conn
	username   string
}

var clients []Client

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
		go handleClientLogin(conn)
	}
}

func handleClientLogin(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		conn.Write([]byte("Please enter your Username! Pick a length between 3-20 characters \r \n"))
		nameinput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		nameinput = strings.TrimSpace(nameinput)
		if len(nameinput) >= 3 && len(nameinput) <= 20 {
			newclient := Client{
				connection: conn,
				username:   nameinput,
			}

			clients = append(clients, newclient)
			go handleConnection(newclient)
			break
		} else {
			conn.Write(([]byte("Username has to be longer >= 3 chars and <= 20 chars \r \n")))
		}
	}

}

func handleConnection(client Client) {
	reader := bufio.NewReader(client.connection)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			removeClient(client)
			broadcast(client, " has left the chat\r\n")
			return
		}
		broadcast(client, msg)
	}
}

func removeClient(client Client) {
	client.connection.Close()
	for i, currentClient := range clients {
		if currentClient == client {
			clients = append(clients[:i], clients[i+1:]...)
		}
	}
}

func broadcast(sender Client, msg string) {
	for _, currentClient := range clients {
		if currentClient != sender {
			currentClient.connection.Write([]byte("[" + currentClient.username + "] " + msg))
		}
	}
}
