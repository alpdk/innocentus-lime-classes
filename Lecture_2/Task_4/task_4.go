/*
	Task: Find minimum value
*/

package main

import (
	"errors"
	"net"
	"os"
	"time"
)

const wait_duration = 1 * time.Second

func main() {
	ip := net.IPv4(13, 61, 191, 208)
	host := net.UDPAddr{IP: ip, Port: 9080}
	conn, err := net.DialUDP("udp", nil, &host)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 100)

	for {
		cur_time := time.Now()

		_, err = conn.Write([]byte("Hello"))
		if err != nil {
			panic(err)
		}

		err := conn.SetReadDeadline(cur_time.Add(wait_duration))
		if err != nil {
			panic(err)
		}

		le, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				println("BOOOOORIIIIING!!!!")
				continue
			} else {
				panic(err)
			}

		}

		print(string(buffer[:le]))
		break
	}
}
