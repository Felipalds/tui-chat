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

	fmt.Println(os.Args)
	address := os.Args[1]
	port := os.Args[2]

	conn, err = net.Dial("tcp", address+":"+port)
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
			encrypted, _ := encryption.Encrypt([]byte(text), aesKey)

			if text == "" {
				continue
			}
			fmt.Fprintf(conn, "%s\n", encrypted)
		}

		if parts[0] == "/msg" {
			text := "ENVIAR_MENSAGEM " + parts[1]
			for i, p := range parts {
				if i > 1 {
					text += " " + p
				}
			}
			encrypted, _ := encryption.Encrypt([]byte(text), aesKey)

			if text == "" {
				continue
			}
			fmt.Fprintf(conn, "%s\n", encrypted)
		}

		if parts[0] == "/list" {
			text := "LISTAR_SALAS"
			encrypted, _ := encryption.Encrypt([]byte(text), aesKey)
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
			encrypted, _ := encryption.Encrypt([]byte(text), aesKey)
			if text == "" {
				continue
			}
			fmt.Fprintf(conn, "%s\n", encrypted)
		}

		if parts[0] == "/ban" {
			text := "BANIR_USUARIO " + parts[1] + " " + parts[2]
			encrypted, _ := encryption.Encrypt([]byte(text), aesKey)
			if text == "" {
				continue
			}
			fmt.Fprintf(conn, "%s\n", encrypted)
		}

		if parts[0] == "/leave" {
			text := "SAIR_SALA " + parts[1]
			encrypted, _ := encryption.Encrypt([]byte(text), aesKey)
			if text == "" {
				continue
			}
			fmt.Fprintf(conn, "%s\n", encrypted)
		}

		if parts[0] == "/close" {
			text := "FECHAR_SALA " + parts[1]
			encrypted, _ := encryption.Encrypt([]byte(text), aesKey)
			if text == "" {
				continue
			}
			fmt.Fprintf(conn, "%s\n", encrypted)
		}

		if parts[0] == "/exit" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from input:", err)
	}
}

func treatMessageFromServer(msg string, conn net.Conn) {
	msg = strings.ReplaceAll(msg, "\n", "")
	if len(msg) <= 2 {
		return
	}
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
		fmt.Println("Receiving a message from server: ", msg)

		treatMessageFromServer(msg, conn)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}
