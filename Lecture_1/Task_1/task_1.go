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

	buffer := make([]byte, 100)
	conn.Write([]byte("Hello"))
	le, _ := conn.Read(buffer)

	print(string(buffer[:le]))
}
