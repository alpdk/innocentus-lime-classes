/*
	Task: make proper protocol implementation
*/

package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path"
	"slices"
)

const packet_size = 16
const start_pay_load = 0x0040
const query_flag = 0x80

const flag_dead_err = 0x00
const flag_info = 0x40
const flag_success = 0x20
const flag_redirect = 0x10
const flag_retry = 0x08
const flag_404 = 0x04
const flag_fragmented = 0x02

const payload_start = 5
const query_size = 7

const data_root = 0xFFFF
const name_root = 0xFFFE

var N = uint16(0)

func put_request_id(send_buffer []byte, req_id uint16) {
	binary.BigEndian.PutUint16(send_buffer[2:4], req_id)
}

func set_request(buffer []byte, id []byte, req_id uint16, resorce_id uint16) {
	// set id
	copy(buffer[0:1], id)

	// set N
	put_request_id(buffer, req_id)

	// set flag
	buffer[4] = query_flag

	// set start payload
	binary.BigEndian.PutUint16(buffer[payload_start:query_size], resorce_id)
}

func recur_protocol_parse(conn *net.UDPConn, send_buffer []byte, receive_buffer []byte, file io.Writer) {
	for {
		le, err := conn.Read(receive_buffer)
		if err != nil {
			panic(err)
		}

		flag := receive_buffer[4]
		//fmt.Printf("%x\n", flag)

		switch {
		case flag == flag_dead_err:
			panic("ABSOLUTE DARK SHIT")
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
				set_request(send_buffer, []byte{0, 0}, N, binary.BigEndian.Uint16(new_data))
				_, err = conn.Write(send_buffer[:query_size])
				if err != nil {
					panic(err)
				}
				recur_protocol_parse(conn, send_buffer, receive_buffer, file)
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
			set_request(send_buffer, []byte{0, 0}, N, binary.BigEndian.Uint16(receive_buffer[payload_start:7]))
			_, err = conn.Write(send_buffer[:query_size])
			if err != nil {
				panic(err)
			}
			//println("The princess at different castle...")
		case flag&flag_success == flag_success:
			//println(flag)
			//println("YUPPY!!!")
			//print(string(receive_buffer[payload_start:le]))
			_, err := file.Write(receive_buffer[payload_start:le])
			if err != nil {
				panic(err)
			}
			return
		}
	}
}

func main() {
	ip := net.IPv4(13, 61, 191, 208)
	host := net.UDPAddr{IP: ip, Port: 9091}
	conn, err := net.DialUDP("udp", nil, &host)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	send_buffer := make([]byte, packet_size)
	receive_buffer := make([]byte, packet_size)

	file_name := bytes.Buffer{}

	set_request(send_buffer, []byte{0, 0}, N, name_root)
	_, err = conn.Write(send_buffer[:query_size])
	if err != nil {
		print("ERROR")
		return
	}

	recur_protocol_parse(conn, send_buffer, receive_buffer, &file_name)

	_, file := path.Split(file_name.String())
	println(file)

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	set_request(send_buffer, []byte{0, 0}, N, data_root)
	_, err = conn.Write(send_buffer[:query_size])
	if err != nil {
		print("ERROR")
		return
	}

	recur_protocol_parse(conn, send_buffer, receive_buffer, f)

	err = f.Close()
	if err != nil {
		panic(err)
	}
}
