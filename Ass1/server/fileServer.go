package main

import "net"
import "fmt"
import "bufio"
import "io/ioutil"

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
		if len(readMessage) >= 4 && readMessage[0:4] == "read" {
			go readFunction(&conn, readMessage)
		}
		fmt.Print("Message is ", string(readMessage))
		// conn.Write([]byte("Able to write \n"))
	}
	conn.Close()
}

func readFunction(conn *net.Conn, readMessage string) {
	fmt.Println("In Read \n")
	fileName := readMessage[5 : len(readMessage)-2]
	content, err := ioutil.ReadFile("./" + fileName)
	fmt.Println("./" + fileName)
	if err == nil {
		(*conn).Write(content)
		// fmt.Print(string(content))
	} else {
		(*conn).Write([]byte("File not found \n"))
		// fmt.Println("File not found")
	}
}
