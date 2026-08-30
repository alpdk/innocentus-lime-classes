/*

Task: "I am from some cool place (TCP EDITION)"

*/

package main

import (
	"net"
	"strconv"
)

func main() {
	_, add, err := net.LookupSRV("text", "tcp", "course.innolime.dev")

	if err != nil {
		panic("LookupSRV error!!!")
	}

	conn, err := net.Dial("tcp", add[0].Target+":"+strconv.Itoa(int(add[0].Port)))

	if err != nil {
		print("Dial error!!!")
		return
	}

	defer conn.Close()

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
