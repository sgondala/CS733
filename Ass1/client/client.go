package main

import (
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
	conn.Write([]byte("Hello\r\n"))
	// content := []byte("test\r\n")
	// time.Sleep(1 * time.Second)
	// conn.Write([]byte("write test 4 4\r\n"))
	// conn.Write(content)
}
