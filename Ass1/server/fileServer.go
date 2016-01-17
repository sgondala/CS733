package main

import "net"
import "fmt"
import "bufio"

// import "io"

// import "strings" // only needed below for sample processing

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
		readMessage, err := bufio.NewReader(conn).ReadString(byte('\n')) // Line is showed as
		if err != nil {
			break
		}
		fmt.Print("Message is ", string(readMessage))
	}
	conn.Close()
}
