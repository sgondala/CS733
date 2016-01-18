package main

import "fmt"
import "net"
import _ "strings"
import "io/ioutil"

func main() {
	fileName := "a"
	content, err := ioutil.ReadFile("./" + fileName)
	if err == nil {
		fmt.Print(string(content))
	} else {
		fmt.Println("File not found")
	}
	// fmt.Println(len(content))
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {

}
