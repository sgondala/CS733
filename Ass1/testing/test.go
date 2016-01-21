package main

import "fmt"
import "net"

// import "strings"
import _ "io/ioutil"
import "time"

var isWrite bool = false

func main() {
	now := time.Now()
	secs := now.Unix()

	fmt.Println(secs)
	i := 1
	fmt.Println(float64(i))
	// fmt.Println(isWrite)
	// fileName := "a"
	// content, err := ioutil.ReadFile("./" + fileName)
	// if err == nil {
	// 	fmt.Print(string(content))
	// } else {
	// 	fmt.Println("File not found")
	// }
	// s := "HelloWorld"
	// words := strings.Split(s, " ")
	// for _, word := range words {
	// 	fmt.Println(word)
	// }
	// fmt.Println(len(content))
}

func tempFunction() {
	fmt.Println("cane")
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {

}
