package main

import "net"
import "fmt"
import "bufio"
import "io/ioutil"
import "strconv"
import "os"
import "strings" // only needed below for sample processing
// import "io"

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
	reader := bufio.NewReader(conn)
	for {
		readMessage, err := reader.ReadString(byte('\n')) // Line is showed as
		if err != nil {
			break
		}
		words := strings.Split(readMessage, " ")

		if words[0] == "write" {
			fileName := words[1]
			numBytes, _ := strconv.Atoi(words[2][:len(words[2])-2])
			contentBytes := make([]byte, numBytes+2)
			for i := 0; i < numBytes+2; i++ {
				contentBytes[i], _ = reader.ReadByte()
			}
			contentBytes = contentBytes[:len(contentBytes)-2]
			fmt.Println(fileName)
			fmt.Println(numBytes)
			fmt.Println(len(contentBytes))
			fmt.Println(string(contentBytes))
			go writeFunction(conn, fileName, numBytes, string(contentBytes))

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
	fileName := readMessage[5 : len(readMessage)-2]
	content, err := ioutil.ReadFile("./" + fileName)
	fmt.Println("./" + fileName)
	if err == nil {
		conn.Write(content)
	} else {
		conn.Write([]byte("File not found \n"))
	}
}
