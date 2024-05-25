package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Welcome to the chat. /help to get the manual")

	go getMessagesFromServer(conn)

	// TODO: study the scanner
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		fmt.Fprintf(conn, "%s\n", text)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from input:", err)
	}
}

func getMessagesFromServer(conn net.Conn) {

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf(msg)
	}
}
