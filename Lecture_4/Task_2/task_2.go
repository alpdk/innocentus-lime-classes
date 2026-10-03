package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8090")

	if err != nil {
		log.Fatal("Error listening:", err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("Error accepting conn:", err)
		}

		go sendFile(conn)
	}
}

func sendFile(conn net.Conn) {
	defer conn.Close()

	str := []byte("Hello, World!")

	for _, byte_data := range str {
		_, err := conn.Write([]byte{byte_data})

		if err != nil {
			log.Printf("Write error: %v\n", err)
			break
		}
	}
}
