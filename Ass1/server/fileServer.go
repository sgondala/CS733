package main

import "net"
import "fmt"
import "bufio"
import "io/ioutil"
import "strconv"
import "os"
import "strings" // only needed below for sample processing

var filesTillNow int64 = 0

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
		words := strings.Split(readMessage, " ")

		if words[0] == "write" {
			fileName := words[1]
			numBytes, _ := strconv.Atoi(words[2])
			contentBytes, _ := bufio.NewReader(conn).ReadString(byte('\n')) //TODO - Assuming that next line always exists
			contentBytes = contentBytes[0 : len(contentBytes)-2]
			go writeFunction(conn, fileName, numBytes, contentBytes)

		} else if words[0] == "delete" {
			fileName := words[1]
			go deleteFunction(conn, fileName[:len(fileName)-2])
		}

		if len(readMessage) >= 4 && readMessage[0:4] == "read" {
			go readFunction(conn, readMessage)
		}
	}
	conn.Close()
}

func deleteFunction(conn net.Conn, fileName string) {
	err := os.Remove(fileName)
	if err == nil {
		conn.Write([]byte("OK\r\n"))
	} else {
		conn.Write([]byte("Error in deletion \n")) // TODO Should check if this is the error
	}
}

func writeFunction(conn net.Conn, fileName string, numBytes int, contentBytes string) {
	conn.Write([]byte("OK " + strconv.FormatInt(filesTillNow, 10) + "\r\n"))
	filesTillNow++ //TODO - Should use concurrency stuff
}

func readFunction(conn net.Conn, readMessage string) {
	fmt.Println("In Read \n")
	fileName := readMessage[5 : len(readMessage)-2]
	content, err := ioutil.ReadFile("./" + fileName)
	fmt.Println("./" + fileName)
	if err == nil {
		conn.Write(content)
		// fmt.Print(string(content))
	} else {
		conn.Write([]byte("File not found \n"))
		// fmt.Println("File not found")
	}
}
