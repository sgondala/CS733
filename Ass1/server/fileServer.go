package main

import "net"
import "fmt"
import "bufio"
import "io/ioutil"
import "strconv"
import "os"
import "strings"
import "sync"
import "time"

var fileVersionNext int64 = 1
var fileVersionMap map[string]int64 = map[string]int64{}
var mutexLock = &sync.Mutex{}

/*
	TODO :- Correct error codes, Read output format, Tests, enable support for expiry
*/

func main() {
	serverMain()
}

func serverMain() {
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
		readMessage, err := reader.ReadString(byte('\n'))
		if err != nil {
			break
		}
		words := strings.Split(readMessage, " ")

		if words[0] == "write" {
			fileName := words[1]
			var numBytes int
			var expTime float64
			if len(words) == 3 {
				numBytes, _ = strconv.Atoi(words[2][:len(words[2])-2])
				expTime = -1
			} else {
				numBytes, _ = strconv.Atoi(words[2])
				expTime, _ = strconv.ParseFloat(words[3][:len(words[3])-2], 64)
			}
			contentBytes := make([]byte, numBytes+2)
			for i := 0; i < numBytes+2; i++ {
				contentBytes[i], _ = reader.ReadByte()
			}
			contentBytes = contentBytes[:len(contentBytes)-2]
			go writeFunction(conn, fileName, numBytes, contentBytes, expTime)

		} else if words[0] == "delete" {
			fileName := words[1]
			go deleteFunction(conn, fileName[:len(fileName)-2])

		} else if words[0] == "cas" {
			fileName := words[1]
			version, _ := strconv.ParseInt(words[2], 10, 64)
			var numBytes int
			var expTime float64
			if len(words) == 4 {
				numBytes, _ = strconv.Atoi(words[3][:len(words[3])-2])
				expTime = -1
			} else {
				numBytes, _ = strconv.Atoi(words[3])
				expTime, _ = strconv.ParseFloat(words[4][:len(words[4])-2], 64)
			}
			contentBytes := make([]byte, numBytes+2)
			for i := 0; i < numBytes+2; i++ {
				contentBytes[i], _ = reader.ReadByte()
			}
			contentBytes = contentBytes[:len(contentBytes)-2]
			go casFunction(conn, fileName, version, numBytes, contentBytes, expTime)

		} else if len(readMessage) >= 4 && readMessage[0:4] == "read" {
			go readFunction(conn, readMessage)
		}
	}
	conn.Close()
}

func deleteFunction(conn net.Conn, fileName string) {
	mutexLock.Lock()
	err := os.Remove(fileName)
	if err == nil {
		delete(fileVersionMap, fileName)
		conn.Write([]byte("OK\r\n"))
	} else {
		conn.Write([]byte("Error in deletion \n"))
	}
	mutexLock.Unlock()
}

func writeFunction(conn net.Conn, fileName string, numBytes int, contentBytes []byte, expTime float64) {
	mutexLock.Lock()
	if !Exists(fileName) {
		fmt.Println("Createad file")
		os.Create(fileName)
	}
	ioutil.WriteFile(fileName, contentBytes, 0644)
	conn.Write([]byte("OK " + strconv.FormatInt(fileVersionNext, 10) + "\r\n"))
	fileVersionNext++
	fileVersionMap[fileName] = fileVersionNext
	if expTime != -1 {
		go deleteFileVersion(fileName, fileVersionNext, expTime)
	}
	mutexLock.Unlock()
}

func deleteFileVersion(fileName string, fileVersion int64, expTime float64) {
	time.Sleep(time.Duration(expTime) * time.Second)
	mutexLock.Lock()
	if fileVersionMap[fileName] == fileVersion {
		err := os.Remove(fileName)
		if err == nil {
			delete(fileVersionMap, fileName)
		}
	} else {
		// fmt.Println("Couldn't delete, file changed")
	}
	mutexLock.Unlock()
}

func readFunction(conn net.Conn, readMessage string) {
	mutexLock.Lock()
	fileName := readMessage[5 : len(readMessage)-2]
	content, err := ioutil.ReadFile("./" + fileName)
	if err == nil {
		conn.Write(content)
	} else {
		conn.Write([]byte("File not found \n"))
	}
	mutexLock.Unlock()
}

func casFunction(conn net.Conn, fileName string, version int64,
	numBytes int, contentBytes []byte, expTime float64) {
	mutexLock.Lock()
	if !Exists(fileName) {
		conn.Write([]byte("File not found \n"))
	} else if fileVersionMap[fileName] == version {
		ioutil.WriteFile(fileName, contentBytes, 0644)
		conn.Write([]byte("OK " + strconv.FormatInt(fileVersionNext, 10) + "\r\n"))
		fileVersionNext++
		fileVersionMap[fileName] = fileVersionNext
		if expTime != -1 {
			go deleteFileVersion(fileName, fileVersionNext, expTime)
		}
	} else {
		conn.Write([]byte("File version mismatch \n"))
	}
	mutexLock.Unlock()
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
