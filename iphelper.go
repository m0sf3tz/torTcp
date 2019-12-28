package main

import "strings"
import "strconv"
import "log"
import "net"
import "fmt"

const PORT = "8082"

type Message struct {
	ClientIP string
	Hash     [8]byte
	FileName string
}

type Response struct {
	FileSize int
	Hash     []byte
}

type Seeders struct {
	SeedIp   string
	Busy     bool
	FileSize int
	Hash     []byte
}

type IpInfo struct {
	CIDR      string
	Broadcast string
	MyIp      string
}

func RemovePort(ip string) string {
	return strings.Split(ip, ":")[0]
}

func IsIpV4(ip string) bool {
	ipSep := strings.Split(ip, "/")
	ipNoCIDR := net.ParseIP(ipSep[0])
	if net.IP.To4(ipNoCIDR) != nil {
		return true
	}
	return false
}

func wAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("failed to convert?")
	}
	return n
}

func GetBroadcastBits(cidr string) string {
	subString := strings.Split(cidr, "/")
	ipString := strings.Split(subString[0], ".")

	var val uint32 = uint32(wAtoi(ipString[0])<<24 + wAtoi(ipString[1])<<16 +
		wAtoi(ipString[2])<<8 + wAtoi(ipString[3]))

	var hostBits uint32 = uint32((1<<(32-wAtoi(subString[1])) - 1))
	var ip32 uint32 = val | hostBits
	var s string

	for i := 0; i < 4; i++ {
		cur := (ip32 >> (8 * i) & 0xFF)
		s = strconv.Itoa(int(cur)) + s
		if i != 3 {
			s = "." + s
		}
	}
	return s
}

//test functions

func GetIps(dev string) IpInfo {

	var interfaceIps []IpInfo

	i, err := net.InterfaceByName(dev)
	if err != nil {
		fmt.Println("failed to find interface, exiting")
		panic(err)
	}

	ipArr, err := i.Addrs()
	for _, ip := range ipArr {
		if IsIpV4(ip.String()) {
			interfaceIps = append(interfaceIps, IpInfo{ip.String(), GetBroadcastBits(ip.String()), strings.Split(ip.String(), "/")[0]})
		}
	}

	if len(interfaceIps) > 1 || (interfaceIps == nil) {
		fmt.Println("This program only supports single-homed interfaces")
		panic(0)
	}
	return interfaceIps[0]
}
