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
var fileExpTimeMap map[string]float64 = map[string]float64{}

func main() {
	serverMain()
}

func serverMain() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			conn.Close()
		} else {
			go singleConnection(conn)
		}
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

		} else if words[0] == "read" {
			fileName := words[1]
			go readFunction(conn, fileName[:len(fileName)-2])

		} else {
			conn.Write([]byte("ERR_CMD_ERR\r\n"))
			// break
		}
	}
	conn.Close()
}

func deleteFunction(conn net.Conn, fileName string) {
	mutexLock.Lock()
	if !Exists(fileName) {
		conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
	} else {
		err := os.Remove(fileName)
		if err == nil {
			delete(fileVersionMap, fileName)
			delete(fileExpTimeMap, fileName)
			conn.Write([]byte("OK\r\n"))
		} else {
			conn.Write([]byte("ERR_INTERNAL\r\n"))
		}
	}
	mutexLock.Unlock()
}

func writeFunction(conn net.Conn, fileName string, numBytes int, contentBytes []byte, expTime float64) {
	mutexLock.Lock()
	if !Exists(fileName) {
		os.Create(fileName)
	}
	ioutil.WriteFile(fileName, contentBytes, 0644)
	fileVersionNext++
	conn.Write([]byte("OK " + strconv.FormatInt(fileVersionNext, 10) + "\r\n"))
	fileVersionMap[fileName] = fileVersionNext
	fileExpTimeMap[fileName] = 0
	if expTime != -1 {
		i := timeInSecsNow()
		fileExpTimeMap[fileName] = float64(i) + expTime
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
			delete(fileExpTimeMap, fileName)
		}
	} else {
		// fmt.Println("Couldn't delete, file changed")
	}
	mutexLock.Unlock()
}

func readFunction(conn net.Conn, fileName string) {
	mutexLock.Lock()
	if !Exists(fileName) {
		conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
	} else {
		content, err := ioutil.ReadFile("./" + fileName)
		if err == nil {
			expTime := fileExpTimeMap[fileName]
			var timeLeft float64
			if expTime == 0 {
				timeLeft = 0
			} else {
				timeLeft = expTime - timeInSecsNow()
				if timeLeft < 0 {
					timeLeft = 0
				}
			}
			conn.Write([]byte("CONTENTS " + strconv.FormatInt(fileVersionMap[fileName], 10) +
				" " + strconv.Itoa(len(content)) + " " + strconv.FormatFloat(timeLeft, 'f', -1, 64) + " " + "\r\n"))
			conn.Write(content)
			conn.Write([]byte("\r\n"))
		} else {
			conn.Write([]byte("ERR_INTERNAL\r\n"))
		}
	}
	mutexLock.Unlock()
}

func casFunction(conn net.Conn, fileName string, version int64,
	numBytes int, contentBytes []byte, expTime float64) {
	mutexLock.Lock()
	if !Exists(fileName) {
		conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
	} else if fileVersionMap[fileName] == version {
		ioutil.WriteFile(fileName, contentBytes, 0644)
		fileVersionNext++
		conn.Write([]byte("OK " + strconv.FormatInt(fileVersionNext, 10) + "\r\n"))
		fileVersionMap[fileName] = fileVersionNext
		fileExpTimeMap[fileName] = 0
		if expTime != -1 {
			fileExpTimeMap[fileName] = timeInSecsNow() + expTime
			go deleteFileVersion(fileName, fileVersionNext, expTime)
		}
	} else {
		currentVersion := fileVersionMap[fileName]
		conn.Write([]byte("ERR_VERSION " + strconv.FormatInt(currentVersion, 10) + "\r\n"))
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

func timeInSecsNow() float64 {
	return float64(time.Now().Unix())
}
