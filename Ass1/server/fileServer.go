package main

import "net"
import "fmt"
import "bufio"
import _ "strings" // only needed below for sample processing

func main() {
	fmt.Println("Launching server...") // listen on all interfaces
	ln, _ := net.Listen("tcp", "localhost:8080")
	for {
		conn, _ := ln.Accept()
		go singleConnection(conn)
	}
}

func singleConnection(conn net.Conn) {
	for {
		readMessage, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message is ", string(readMessage))
	}
}
