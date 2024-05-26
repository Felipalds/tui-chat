package main

import (
	"bufio"
	"fmt"
	"github.com/Felipalds/tui-chat.git/encryption"
	"net"
	"os"
	"strings"
)

var conn net.Conn

func main() {

	var err error
	conn, err = net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Welcome to the chat. /help to get the manual")

	go getMessagesFromServer(conn)

	// TODO: study the bufio lib
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

func treatMessageFromServer(msg string, conn net.Conn) {

	buffParts := strings.Split(msg, " ")
	buffParts[len(buffParts)-1] = strings.ReplaceAll(buffParts[len(buffParts)-1], "\x00", "")

	requestType := strings.ToUpper(buffParts[0])

	if requestType == "CHAVE_PUBLICA" {
		fmt.Println(requestType[1])
		fmt.Println(len(buffParts[1]))
		aesKey, err := encryption.ReadPublicKey(buffParts[1])
		if err != nil {
			panic(err)
		}
		fmt.Println("len aes key", len(aesKey))
		fmt.Fprintf(conn, "CHAVE_SIMETRICA "+aesKey+"\n")
	}

}

func getMessagesFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')

		treatMessageFromServer(msg, conn)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}
