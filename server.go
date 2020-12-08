package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var servchat []string

var (
	//I store my connections in a map
	openConnections = make(map[net.Conn]bool)
	newConnection   = make(chan net.Conn)
	deadConnection  = make(chan net.Conn)
)

func showChat() {
	for _, mes := range servchat {
		fmt.Println(mes)
	}
}

func logServerChat() {
	file, err := os.Create("Server.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for v := range servchat {
		file.WriteString(servchat[v])

	}

}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func broadcastMessage(conn net.Conn) {

	for {

		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		fmt.Println(message)
		fmt.Println("-----------------------------------------------")
		servchat = (append(servchat, message))

		if err != nil {
			break
		}

		for item := range openConnections {
			if item != conn {
				item.Write([]byte(message))

			}

		}

	}
	deadConnection <- conn
}

func main() {
	fmt.Println("Starting...")
	time.Sleep(time.Second * 2)
	fmt.Println("Running")
	ln, err := net.Listen("tcp", ":8080")
	logFatal(err)

	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			logFatal(err)

			openConnections[conn] = true
			newConnection <- conn

		}
	}()

	for {
		select {
		case conn := <-newConnection:
			go broadcastMessage(conn)

		case conn := <-deadConnection:
			for item := range openConnections {
				if item == conn {
					break
				}
			}
			delete(openConnections, conn)
			var opc string

			for opc != "3" {
				fmt.Println("1)Show chat history\n2)Back up chat\n3)continue\n0)Exit")
				fmt.Scanln(&opc)

				switch opc {
				case "1":
					showChat()
				case "2":
					logServerChat()
				case "3":
					fmt.Println("continue")
				case "0":
					fmt.Println("Terminated")
					ln.Close()

				}
			}

		}

	}

}
