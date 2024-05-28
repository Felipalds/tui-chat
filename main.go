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
var auth bool

func main() {

	var err error
	auth = false
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
		if parts[0] == "/login" {
			text := "AUTENTICACAO " + parts[1]
			fmt.Println(text)
			fmt.Fprintf(conn, "%s\n", text)
		}
		if parts[0] == "/signup" {
			text := "REGISTRO " + parts[1]
			fmt.Fprintf(conn, "%s\n", text)
		}
		if parts[0] == "/create" {
			text := "CRIAR_SALA " + parts[1] + " " + parts[2]
			if len(parts) == 4 {
				text += " " + parts[3]
			}
			encrypted, _ := encryption.Encrypt(text, aesKey)

			if text == "" {
				continue
			}
			fmt.Fprintf(conn, "%s\n", encrypted)
		}

		if parts[0] == "/msg" {
			text := "ENVIAR_MENSAGEM " + parts[1]
			for _, p := range parts {
				text += " " + p
			}
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

		if parts[0] == "/join" {
			text := "ENTRAR_SALA " + parts[1]
			if len(parts) == 3 {
				text += " " + parts[2]
			}
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
	msg = strings.ReplaceAll(msg, "\n", "")
	var err error

	if auth {
		msg, err = encryption.Decrypt(msg, aesKey)
		if err != nil {
			fmt.Println("Error decrypting message:", err)
		}
	}

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
		auth = true
	} else {
		for _, p := range buffParts {
			fmt.Print(p, " ")
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
