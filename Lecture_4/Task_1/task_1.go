package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
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

	f, err := os.Open("/home/alpdk/gitRepos/innocentus-lime-classes/Lecture_4/images.jpeg")

	if err != nil {
		log.Fatal("File reading error:", err)
	}

	defer f.Close()

	buffer := make([]byte, 8)

	for {
		n, err := f.Read(buffer)

		if err != nil {
			if errors.Is(err, io.EOF) {
				_, err = conn.Write(buffer[:n])
				break
			}

			log.Fatal("Read error:", err)
		}

		_, err_wr := conn.Write(buffer[:n])

		if err_wr != nil {
			log.Fatal("Error buffer sending:", err_wr)
		}
	}
}
