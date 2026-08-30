/*

Task: "I am from some cool place (UDP EDITION)"

*/

package main

import (
	"net"
	"strconv"
)

func main() {
	_, add, err := net.LookupSRV("text", "udp", "course.innolime.dev")

	if err != nil {
		panic("LookupSRV error!!!")
	}

	conn, err := net.Dial("udp", add[0].Target+":"+strconv.Itoa(int(add[0].Port)))

	if err != nil {
		print("Dial error!!!")
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte("file"))

	if err != nil {
		print("Write error!!!")
		return
	}

	buffer := make([]byte, 7)

	for {
		le, err := conn.Read(buffer)

		if err != nil {
			print("Read error!!!")
			return
		}

		print(string(buffer[:le]))
	}
}
