package main

import "os"
import "fmt"
import "net"
import "log"
import "encoding/json"

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func check(err error) {
	if err != nil {
		panic(0)
	}
}

func Stat(name string) Response {
	fileInfo, err := os.Stat(name)
	check(err)

	var ret Response
	ret.FileSize = int(fileInfo.Size())
	ret.Hash = sha256sum(name)
	return (ret)
}

func respond(m Message) {
	//get the IP of the requestee
	ipRequestee := m.ClientIP + ":8081"

	if !Exists(m.FileName) {
		return
	}

	conn, err := net.Dial("udp", ipRequestee)
	if err != nil {
		fmt.Printf("tried to connect to %s \n", ipRequestee)
		log.Fatal("Could not resond to ClientIP")
	}

	response := Stat(m.FileName)
	b, err := json.Marshal(response)
	check(err)

	conn.Write(b)

	conn.Close()

}

func leecherListener() {
	b := make([]byte, 256)

	for {
		ln, err := net.ListenPacket("udp", ":8081")
		if err != nil {
			log.Fatal(err)
		}

		n, addr, err := ln.ReadFrom(b)
		if err != nil {
			log.Fatal("Could accept")
		}

		fmt.Printf("Incomming from %s \n ", addr.String())

		var m Message

		var er error = json.Unmarshal(b[0:n], &m)
		if er != nil {
			fmt.Println("Could not unmarshal Message")
		}

		respond(m)

		fmt.Println(m)

		ln.Close()
	}

}
