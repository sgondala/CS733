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
	conn.Write([]byte("read a\r\n"))
	conn.Write([]byte("read b\r\n"))
	// for {
	// 	readMessage, err := bufio.NewReader(conn).ReadString(byte('\n')) // Line is showed as
	// 	if err != nil {
	// 		break
	// 	}
	// 	fmt.Println(string(readMessage))
	// 	// conn.Write([]byte("Able to write \n"))
	// }
}
