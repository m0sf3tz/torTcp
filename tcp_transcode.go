package main

import "fmt"
import "net"
import "io"

const SIZE_HEADER_REQUEST int = 4
const CHUNK_SIZE int = 1 << 20

func tcpTranser(conn net.Conn, size int, chunk []byte) {
	i := 0

	for i < size {
		if i+CHUNK_SIZE > size {
			n, err := conn.Write(chunk[i:])
			if err != nil {
				fmt.Println(err)
				panic(0)
			}
			i = i + n
		} else {
			n, err := conn.Write(chunk[i : i+CHUNK_SIZE])
			if err != nil {
				fmt.Println(err)
				panic(0)
			}
			i = i + n
		}
	}
}

func tcpRecieve(conn net.Conn, size int) []byte {
	recv := make([]byte, 1024)
	b := make([]byte, 1<<20)

	i := 0
	for i < size {
		n, err := conn.Read(recv)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				panic(0)
			}
			//socket is closed by other end, return what we have
			return b
		}
		b = append(b[0:i], recv[0:n]...)
		i = i + n
	}
	return b
}

//byte 0-3 [requested chunk]
//byte 4   [len of File name]
//byte 4-n [filename]
func encodeRequest(fileName string, chunk int) []byte {
	ret := make([]byte, 4)

	for i := 0; i < 4; i++ {
		ret[i] = byte(chunk & 0xFF)
		chunk = chunk >> (8)
	}

	ba := []byte(fileName)
	bs := byte(len(ba))

	ret = append(ret, bs)
	ret = append(ret, ba[:]...)

	return ret
}

func decodeRequest(e []byte) (string, int) {
	var chunk int

	for i := 0; i < 4; i++ {
		chunk = chunk + int(e[i])<<(i*8)
	}

	bs := int(e[4])

	fileName := string(e[5 : 5+bs])

	return fileName, chunk
}

func encodeResponse(chunk int, size int, hash [32]byte) []byte {
	ret := make([]byte, 36)

	for i := 0; i < 4; i++ {
		ret[i] = byte(chunk & 0xFF)
		chunk = chunk >> (8)
	}

	for i := 0; i < 4; i++ {
		ret[i+4] = byte(size & 0xFF)
		size = size >> (8)
	}

	ret = append(ret[0:8], hash[:]...)
	return ret
}

func decodeResponse(e []byte) (int, int, [32]byte) {
	var chunk int
	var size int
	var hash [32]byte

	for i := 0; i < 4; i++ {
		chunk = chunk + int(e[i])<<(i*8)
	}
	for i := 0; i < 4; i++ {
		size = size + int(e[i+4])<<(i*8)
	}

	copy(hash[:], e[8:])

	return chunk, size, hash
}
