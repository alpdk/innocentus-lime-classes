/*
	Task: make proper protocol implementation
*/

package main

import (
	"encoding/binary"
	"net"
	"slices"
)

const packet_size = 16
const start_pay_load = 0x0040
const query_flag = 0x80

const flag_info = 0x40
const flag_success = 0x20
const flag_redirect = 0x10
const flag_retry = 0x08
const flag_404 = 0x04
const flag_fragmented = 0x02

const payload_start = 5
const query_size = 7

var N = uint16(0)

func put_request_id(send_buffer []byte, req_id uint16) {
	binary.BigEndian.PutUint16(send_buffer[2:4], req_id)
}

func recur_protocol_parse(conn *net.UDPConn, send_buffer []byte, receive_buffer []byte) {
	for {
		le, err := conn.Read(receive_buffer)
		if err != nil {
			panic(err)
		}

		flag := receive_buffer[4]
		//fmt.Printf("%x\n", flag)

		switch {
		case flag&flag_404 == flag_404:
			//println("S.H.I.T. XP")
			panic("404")
		case flag&flag_info == flag_info:
			//println("But actually...")
			//println(string(receive_buffer[payload_start:le]))
		case flag&flag_fragmented == flag_fragmented:
			bytes_of_fragment_ids := slices.Clone(receive_buffer[payload_start:le])
			for new_data := range slices.Chunk(bytes_of_fragment_ids, 2) {
				// set new N
				N += 1
				copy(send_buffer[payload_start:query_size], new_data)
				_, err = conn.Write(send_buffer[:query_size])
				if err != nil {
					panic(err)
				}
				recur_protocol_parse(conn, send_buffer, receive_buffer)
			}
			return
		case flag&flag_retry == flag_retry:
			_, err = conn.Write(send_buffer[:query_size])
			if err != nil {
				panic(err)
			}
		case flag&flag_redirect == flag_redirect:
			// set new N
			N += 1
			put_request_id(send_buffer, N)
			copy(send_buffer[payload_start:query_size], receive_buffer[payload_start:7])
			_, err = conn.Write(send_buffer[:query_size])
			if err != nil {
				panic(err)
			}
			//println("The princess at different castle...")
		case flag&flag_success == flag_success:
			//println(flag)
			//println("YUPPY!!!")
			print(string(receive_buffer[payload_start:le]))
			return
		}
	}
}

func main() {
	ip := net.IPv4(13, 61, 191, 208)
	host := net.UDPAddr{IP: ip, Port: 9090}
	conn, err := net.DialUDP("udp", nil, &host)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	send_buffer := make([]byte, packet_size)

	// set id
	copy(send_buffer[0:1], []byte{0, 0})

	// set N
	put_request_id(send_buffer, N)

	// set N
	send_buffer[4] = query_flag

	// set start payload
	binary.BigEndian.PutUint16(send_buffer[payload_start:query_size], start_pay_load)

	_, err = conn.Write(send_buffer[:query_size])
	if err != nil {
		print("ERROR")
		return
	}

	receive_buffer := make([]byte, packet_size)

	recur_protocol_parse(conn, send_buffer, receive_buffer)
}
