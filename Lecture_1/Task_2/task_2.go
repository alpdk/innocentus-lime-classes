/*
	Task name: Max number
*/

package main

import "net"
import "strconv"
import "math"

const numCount = 10

func main() {
	ip := net.IPv4(13, 61, 191, 208)
	host := net.UDPAddr{IP: ip, Port: 8090}
	conn, err := net.DialUDP("udp", nil, &host)

	if err != nil {
		print("AAAA ERROR!!!")
		return
	}
	defer conn.Close()

	buffer := make([]byte, 100)
	conn.Write([]byte("Numbers"))

	max_value := math.MinInt32

	for range numCount {
		le, err := conn.Read(buffer)
		if err != nil {
			print("AAAA ERROR!!!")
			return
		}

		value, err := strconv.Atoi(string(buffer[:le]))
		if err != nil {
			print("AAAA ERROR!!!")
			return
		}

		max_value = max(max_value, value)
	}

	println(max_value)
}
