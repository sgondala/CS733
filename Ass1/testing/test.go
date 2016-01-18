package main

import "fmt"
import "net"

// import "strings"
import _ "io/ioutil"

var isWrite bool = false

func main() {
	fmt.Println(0)
	go tempFunction()
	fmt.Println(1)
	tempFunction()
	fmt.Println(2)
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
