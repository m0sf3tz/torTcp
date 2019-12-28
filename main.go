package main

import "net"
import "log"
import "encoding/json"
import "fmt"
import "os"
import "time"
import "sync"

var syncro chan int

var al AtomicList

type AtomicList struct {
	mu      sync.Mutex
	seeders []Seeders
}

func Discover(i IpInfo) {

	conn, err := net.Dial("udp", "192.168.0.255:8081")
	if err != nil {
		log.Fatal("Could not dial")
	}

	m := Message{i.MyIp, [8]byte{1, 1, 1, 1}, os.Args[2]}

	b, _ := json.Marshal(m)

	conn.Write(b)

	conn.Close()
}

func listen(i IpInfo) {

	b := make([]byte, 256)

	listenIp := i.MyIp + ":8081"

	ln, err := net.ListenPacket("udp", listenIp)
	if err != nil {
		log.Fatal(err)
	}

	syncro <- 1

	for {

		ln.SetReadDeadline(time.Now().Add(time.Second))

		i, addr, err := ln.ReadFrom(b)
		if err != nil {

			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				// Timeout error

				syncro <- 1
				return
			}

		}

		fmt.Printf("Incomming from %s \n ", addr.String())

		var r Response
		err = json.Unmarshal(b[0:i], &r)
		check(err)

		al.mu.Lock()
		fmt.Println("found a seed!")
		al.seeders = append(al.seeders, Seeders{addr.String(), false, r.FileSize, r.Hash})
		al.mu.Unlock()
	}
}

func seederSanityCheck(seeders []Seeders) {
	fileSizeIndex0 := seeders[0].FileSize

	for _, seed := range seeders {
		if seed.FileSize != fileSizeIndex0 {
			fmt.Println("different file sizes")
			panic(0)
		}
	}
}

func main() {

	if len(os.Args) == 1 {
		go leecherListener() //deals with finding peers
		seedTcpListener()    //deals with files
	}

	if len(os.Args) != 3 {
		fmt.Println("must supply 2 arguments, 1) interface, 2) filename")
		panic(0)
	}

	iFace := os.Args[1]
	fileName := os.Args[2]

	syncro = make(chan int)

	ipInfo := GetIps(iFace)

	go listen(ipInfo)

	<-syncro

	Discover(ipInfo)

	<-syncro

	if len(al.seeders) == 0 {
		fmt.Println("no seeder found, exiting")
		return
	}
	//delete the file if it already exists in our file-system
	deleteFile(fileName)

	//check that all of our seeders agree with each other
	seederSanityCheck(al.seeders)

	//create a worker list
	ca := createChunkArray(al.seeders[0].FileSize, fileName)

	//starts the TCP process to get the files
	fetchTop(ca, len(al.seeders))

}
