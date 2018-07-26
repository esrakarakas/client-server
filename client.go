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

	ip := ipAddress()
	port := portEntry()
	port2 := strconv.Itoa(port)
	ip2 := ip.String()
	var listen string = ip2 + ":" + port2
	connect(listen)

}
func ipAddress() net.IP {
	var ipv4Addr net.IP
	for {
		var ipv4 string
		fmt.Println("Enter the IP address of the server you want to connect: ")
		fmt.Scanln(&ipv4)

		ips, err := net.LookupIP(ipv4)
		if err != nil {
			fmt.Println("Invalid domain name or IP address. ")

		}

		yeni := ips[0].String()
		ipv4Addr = net.ParseIP(yeni)
		if ipv4Addr == nil {
			fmt.Println("the IP address you entered is incorrect.")
			continue
		} else {
			break
		}
		kontrol := ipv4Addr.DefaultMask()
		if kontrol == nil {
			fmt.Println("Invalid IP address. ")
		} else {
			break
		}
	}
	return ipv4Addr

}

func portEntry() int {
	var port string
	var x int
	var err error
	for {
		fmt.Println("Enter the port  you want to connect: ")
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

func connect(listen string) {
	conn, _ := net.Dial("tcp", listen)
	if conn != nil {
		fmt.Println("-----Launching server-----")
	}
	r := bufio.NewReader(conn)
	publicKeyByte := make([]byte, 270)
	_, err := io.ReadFull(r, publicKeyByte)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyByte)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	clientPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	clientPublicKey := &clientPrivateKey.PublicKey
	sent := x509.MarshalPKCS1PublicKey(clientPublicKey)
	conn.Write(sent)

	for {

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

		r := bufio.NewReader(conn)
		message := make([]byte, 256)
		_, errr := io.ReadFull(r, message)

		if errr != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		plainText, err := rsa.DecryptPKCS1v15(rand.Reader, clientPrivateKey, []byte(message))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		m := string(plainText[:])
		fmt.Println("message received:", m)

	}

}
