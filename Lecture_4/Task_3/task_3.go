package main

import (
	"encoding/binary"
	"net"
	"strings"
)

var transaction_ID = 1

func setup_request(buffer []byte) {
	// set transaction ID
	binary.BigEndian.PutUint16(buffer[0:2], uint16(transaction_ID))

	binary.BigEndian.PutUint16(buffer[2:4], 0x0100)

	binary.BigEndian.PutUint16(buffer[4:6], 0x0001)

	binary.BigEndian.PutUint16(buffer[6:8], 0x0000)

	binary.BigEndian.PutUint16(buffer[8:10], 0x0000)

	binary.BigEndian.PutUint16(buffer[10:12], 0x0000)
}

func main() {
	//data, _ := net.LookupIP("coursens.innolime.dev")
	header_buffer := make([]byte, 12)

	setup_request(header_buffer)

	addres_str := strings.Split("coursens.innolime.dev", ".")
	body_buffer := []byte{}

	for _, data_str := range addres_str {
		body_buffer = append(body_buffer, byte(len(data_str)))
		body_buffer = append(body_buffer, data_str...)
	}

	body_buffer = append(body_buffer, 0)

	var record uint16 = 0x0001
	var class_code uint16 = 0x0001

	request_buffer := append(header_buffer, body_buffer...)
	request_buffer = binary.BigEndian.AppendUint16(request_buffer, record)
	request_buffer = binary.BigEndian.AppendUint16(request_buffer, class_code)

	//print(buffer)

	ip := net.IPv4(1, 1, 1, 1)
	host := net.UDPAddr{IP: ip, Port: 53}
	conn, err := net.DialUDP("udp", nil, &host)
	defer conn.Close()

	if err != nil {
		panic(err)
	}

	_, err = conn.Write(request_buffer)

	if err != nil {
		print("ERROR")
		return
	}

	receive_buffer := make([]byte, 256)
	_, err = conn.Read(receive_buffer)

	if err != nil {
		print("ERROR")
		return
	}

	receive_query_buffer := receive_buffer[(len(request_buffer) + 12):]

	for i, value := range receive_query_buffer {
		print(value)

		if i < 3 {
			print(".")
		} else {
			break
		}
	}
}
