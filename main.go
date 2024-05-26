package main

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"github.com/Felipalds/tui-chat.git/encryption"
	"net"
	"os"
	"strings"
)

var conn net.Conn
var aesKey []byte
var pk *rsa.PublicKey

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
		parts := strings.Split(scanner.Text(), " ")
		if parts[0] == "/autenticacao" {
			text := "AUTENTICACAO " + parts[1] + "\n"
			fmt.Fprintf(conn, "%s\n", text)
		}
		if parts[0] == "/registro" {
			text := "REGISTRO " + parts[1] + "\n"
			fmt.Fprintf(conn, "%s\n", text)
		}
		if parts[0] == "/criar-sala" {
			text := "CRIAR_SALA " + parts[1] + " " + parts[2] + " " + parts[3]
			encrypted, _ := encryption.Encrypt(text, aesKey)

			if text == "" {
				continue
			}
			fmt.Fprintf(conn, "%s\n", encrypted)
		}

		if parts[0] == "/list" {
			text := "LISTAR_SALAS"
			encrypted, _ := encryption.Encrypt(text, aesKey)
			if text == "" {
				continue
			}
			fmt.Fprintf(conn, "%s\n", encrypted)
		}

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
		aesKeyE, aesBytes, err := encryption.ReadPublicKey(buffParts[1])
		aesKey = aesBytes
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(conn, "CHAVE_SIMETRICA "+aesKeyE+"\n")
	} else {
		for _, p := range buffParts {
			fmt.Println(p)
		}
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
