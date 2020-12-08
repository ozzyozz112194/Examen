package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var chat []string

func showChat() {
	for _, mes := range chat {
		fmt.Println(mes)
	}
}

func logChat(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for v := range chat {
		file.WriteString(chat[v])

	}

}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func write(connection net.Conn, username string) {
	var counter int = 0
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		if counter == 1 {
			break
		}

		message = fmt.Sprintf("%s:-%s\n", username, strings.Trim(message, "\r\n"))
		connection.Write([]byte(message))
		chat = (append(chat, message))
		counter++

	}
}

func read(connection net.Conn) {

	for {

		reader := bufio.NewReader(connection)
		message, err := reader.ReadString('\n')

		if err == io.EOF {
			connection.Close()
			fmt.Println("Conection closed")
			os.Exit(0)
		}
		fmt.Println(message)
		fmt.Println("-----------------------------------------------")
		chat = (append(chat, message))
	}

}

func main() {

	connection, err := net.Dial("tcp", "localhost:8080")
	logFatal(err)
	go read(connection)

	defer connection.Close()
	fmt.Println("Enter username:")

	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	logFatal(err)

	username = strings.Trim(username, "\r\n")

	welcomeMsg := fmt.Sprintf("Welcom %s, say ho to your friends.\n", username)

	fmt.Println(welcomeMsg)
	var opc string

	for opc != "0" {

		fmt.Println("1)Write\n2)Show chat history\n3)Back up chat\n0)Exit chat")
		fmt.Scanln(&opc)

		switch opc {
		case "1":
			write(connection, username)

		case "2":
			showChat()
		case "3":
			fmt.Println("Created backup")
			var conc string = username + ".txt"
			logChat(conc)
		case "0":
			fmt.Println("Exit")

		}

	}

}
