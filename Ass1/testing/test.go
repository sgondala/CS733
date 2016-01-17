package main

import (
	"fmt"
	"net"
)

func main() {
	l, _ := net.Listen("tcp", "localhost:8080")
	for {
		conn, _ := l.Accept()
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	message, _ := conn.Read(buf)
	fmt.Print("Message is ", string(message), "\n")
}
