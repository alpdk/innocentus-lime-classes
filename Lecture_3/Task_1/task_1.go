/*

Task: Another "Hello World!"

*/

package main

import (
	"net"
	"strconv"
)

func main() {
	_, add, err := net.LookupSRV("hello", "udp", "course.innolime.dev")

	if err != nil {
		panic("LookupSRV error!!!")
	}

	conn, err := net.Dial("udp", add[0].Target+":"+strconv.Itoa(int(add[0].Port)))

	if err != nil {
		print("Dial error!!!")
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte("Hello"))

	if err != nil {
		print("Write error!!!")
		return
	}

	buffer := make([]byte, 100)
	le, err := conn.Read(buffer)

	if err != nil {
		print("AAA ERROR!!!")
		return
	}

	print(string(buffer[:le]))
}
