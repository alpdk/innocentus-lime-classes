/*

Task: "SPOOKY MONTH!!!"

*/

package main

import (
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	_, add, err := net.LookupSRV("spooky", "tcp", "course.innolime.dev")

	if err != nil {
		panic("LookupSRV error!!!")
	}

	conn, err := net.Dial("tcp", add[0].Target+":"+strconv.Itoa(int(add[0].Port)))

	if err != nil {
		print("Dial error!!!")
		return
	}

	defer conn.Close()

	f, err := os.OpenFile("image.jpeg", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 128)

	for {
		le, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				print("THE EOF... THE EOF IS REAL!!!")
				break
			} else {
				panic("Read error!!!")
				return
			}
		}

		_, err = f.Write(buffer[:le])
		if err != nil {
			panic("Write error!!!")
		}
	}

	err = f.Close()
	if err != nil {
		panic("Close error!!!")
	}
}
