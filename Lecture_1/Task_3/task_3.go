/*
	Task name: Protocol
*/

package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"slices"
)

const packet_size = 16
const start_pay_load = 0x0041
const query_flag = 0x80

const flag_info = 0x40
const flag_success = 0x20
const flag_redirect = 0x10
const flag_retry = 0x08
const flag_404 = 0x04
const flag_fragmented = 0x02

const payload_start = 5

func main() {
	ip := net.IPv4(13, 61, 191, 208)
	host := net.UDPAddr{IP: ip, Port: 9000}
	conn, err := net.DialUDP("udp", nil, &host)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	N := uint16(0)

	send_buffer := make([]byte, packet_size)

	// set id
	copy(send_buffer[0:1], []byte{0, 0})

	// set N
	binary.BigEndian.PutUint16(send_buffer[2:4], N)

	// set N
	send_buffer[4] = query_flag

	// set start payload
	binary.BigEndian.PutUint16(send_buffer[5:7], start_pay_load)

	_, err = conn.Write(send_buffer[:7])
	if err != nil {
		print("ERROR")
		return
	}

	receive_buffer := make([]byte, packet_size)

	for range 100 {
		le, err := conn.Read(receive_buffer)
		if err != nil {
			print("ERROR")
			return
		}

		flag := receive_buffer[4]
		fmt.Printf("%x\n", flag)

		switch {
		case flag&flag_404 == flag_404:
			println("S.H.I.T. XP")
			return
		case flag&flag_info == flag_info:
			println("But actually...")
			println(string(receive_buffer[payload_start:]))
		case flag&flag_fragmented == flag_fragmented:
			receive_buffer_copy := make([]byte, packet_size)
			copy(receive_buffer_copy, receive_buffer)
			le_old := le

			for new_data := range slices.Chunk(receive_buffer_copy[payload_start:le_old], 2) {
				copy(send_buffer[payload_start:7], new_data)
				conn.Write(send_buffer[:7])
				le, _ := conn.Read(receive_buffer)

				print(string(receive_buffer[payload_start:le]))
			}
			return
		case flag&flag_retry == flag_retry:
			conn.Write(send_buffer[:7])
		case flag&flag_redirect == flag_redirect:
			// set new N
			N += 1
			binary.BigEndian.PutUint16(send_buffer[2:4], N)
			copy(send_buffer[payload_start:7], receive_buffer[payload_start:7])
			conn.Write(send_buffer[:7])
			println("The princess at different castle...")
		case flag&flag_success == flag_success:
			println(flag)
			println("YUPPY!!!")
			println(string(receive_buffer[payload_start:]))
			return
		}
	}
}
