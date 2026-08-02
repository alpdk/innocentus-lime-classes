/*
	Task: Find minimum value
*/

package main

import "net"
import "strconv"
import "math"

const numCount = 20

func main() {
	ip := net.IPv4(13, 61, 191, 208)
	host := net.UDPAddr{IP: ip, Port: 9000}
	conn, err := net.DialUDP("udp", nil, &host)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 100)
	_, err = conn.Write([]byte("Numbers"))
	if err != nil {
		panic(err)
	}

	min_value := math.MaxInt32

	for range numCount {
		le, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}

		value, err := strconv.Atoi(string(buffer[:le]))
		if err != nil {
			println(string(buffer[:le]))
			continue
		}

		min_value = min(min_value, value)
	}

	println(min_value)
}
