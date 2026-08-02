/*
	Task name: Hello World!
*/

package main

import "net"

func main() {
	ip := net.IPv4(13, 61, 191, 208)
	host := net.UDPAddr{IP: ip, Port: 8080}
	conn, err := net.DialUDP("udp", nil, &host)

	if err != nil {
		print("AAAA ERROR!!!")
		return
	}
	defer conn.Close()

	_, err := conn.Write([]byte("Hello"))
	if err != nil {
		print("AAAA ERROR!!!")
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
