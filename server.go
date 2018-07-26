package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {

	port := portEntry()
	port2 := strconv.Itoa(port)
	port2 = ":" + port2
	listenPort(port2)

}

func portEntry() int {
	var port string
	var x int
	var err error
	for {
		fmt.Println("enter the port you want to open the connection to: ")
		fmt.Scanln(&port)
		x, err = net.LookupPort("tcp", port)
		if err != nil {
			fmt.Println("Invalid port entry. ")
			
		} else {
			break
		}
	}
	return x
}
func listenPort(t string) {
	ln, _ := net.Listen("tcp", t)
	conn, _ := ln.Accept()
	if conn != nil {
		fmt.Println("-----Launching client-----")
	}

	serverPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	serverPublicKey := &serverPrivateKey.PublicKey
	sent := x509.MarshalPKCS1PublicKey(serverPublicKey)
	conn.Write(sent)

	a := bufio.NewReader(conn)
	publicKeyByte := make([]byte, 270)
	_, er := io.ReadFull(a, publicKeyByte)
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyByte)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if er != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		r := bufio.NewReader(conn)
		message := make([]byte, 256)
		_, err := io.ReadFull(r, message)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		plainText, err := rsa.DecryptPKCS1v15(rand.Reader, serverPrivateKey, []byte(message))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		m := string(plainText[:])
		fmt.Println("message received:", m)
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Text to send: ")
		text, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, text)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		conn.Write(ciphertext)

	}
}
