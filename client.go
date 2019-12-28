package main

import "net"

import "fmt"

//import "encoding/hex"
import "crypto/sha256"

func leech(conn net.Conn, fileName string, chunk int) (bool, []byte) {

	rh := encodeRequest(fileName, chunk)

	conn.Write(rh)

	//get back the confirmation chunk number and also the hash
	recv := make([]byte, 1024)
	n, err := conn.Read(recv)
	if err != nil {
		return false, nil
	}

	_, size, hash := decodeResponse(recv[0:n])
	//fmt.Printf("Seed is sending us chunk %d with size %d and hash %s \n", confirmationChunk, size, string(hex.EncodeToString(hash[:])))

	//TODO: add a sanity check for requested chunk vs the one the server sent us

	//sync-write, the server will not write to us until it gets this,
	//otherwise, we might read parts of the message in the previous read (for teh confirmation chunk)

	conn.Write([]byte{0})
	data := tcpRecieve(conn, size)

	hashR := sha256.Sum256(data[:size])
	//hashRs := hex.EncodeToString(hashR[:])

	fmt.Printf("size of data = %d \n", len(data))

	if hashR != hash {
		return false, nil
	}

	fmt.Printf("size of data = %d \n", len(data))

	return true, data
}

func getChunk(fileName string, chunk int, ip string) bool {

	conn, err := net.Dial("tcp", ip)

	check(err)

	defer conn.Close()

	r, b := leech(conn, fileName, chunk)

	if r {
		err := writeData(fileName, chunk, b)
		if err != true {
			fmt.Println("Failed to write")
			panic(0)
		}
	}
	return r
}

func createChunkArray(size int, fileName string) []chunkStruct {
	chunks := size / (CHUNK_SIZE)
	if size%CHUNK_SIZE != 0 {
		chunks++
	}

	var ca []chunkStruct
	for i := 0; i < chunks; i++ {
		ca = append(ca, chunkStruct{i, false, fileName})
	}
	return ca
}
