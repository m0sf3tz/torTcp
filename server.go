package main

import "net"
import "fmt"

func seedTcpWorker(conn net.Conn) bool {

	defer conn.Close()

	header := make([]byte, 1024)

	//we don't know how big the header is, we will read a big chunk and rely on conn.Read
	//to let us know how big the header was (using n)
	n, err := conn.Read(header)
	if err != nil {
		panic(0)
	}

	fileName, chunk := decodeRequest(header[0:n])

	fmt.Printf("Leech is requesting file=%s, chunk=%d \n", fileName, chunk)

	chunkResponse := readData(fileName, chunk)

	//sending back the chunk and also the hash

	_, err = conn.Write(encodeResponse(chunk, chunkResponse.Len, chunkResponse.Hash))

	if err != nil {
		panic(0)
	}

	//syncronization read
	conn.Read(header)

	//time.Sleep(time.Second)
	tcpTranser(conn, chunkResponse.Len, chunkResponse.Bytes)

	return true
}

func seedTcpListener() {

	fmt.Println("starting in seed mode")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println(err)
		panic(0)
	}

	defer ln.Close()

	for {
		// accept connection on port
		conn, _ := ln.Accept()

		go seedTcpWorker(conn)

	}
}
