package main

import (
	_ "bufio"
	"fmt"
	"net"
	_ "time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	if err != nil {
		fmt.Println("Server not up")
	}
	conn.Write([]byte("write a 10\r\n1234\r\n7890\r\n"))
	// conn.Write([]byte(""))
}
