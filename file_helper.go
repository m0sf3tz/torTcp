package main

import "os"
import "log"
import "crypto/sha256"
import "io"
import "fmt"

import "encoding/hex"

func HeaderPrintOld(s []SliceOld) {
	for _, v := range s {
		fmt.Printf("nth    := %d \n", v.Nth)
		fmt.Printf("offset := %d \n", v.Offset)
		fmt.Printf("hash   := %s \n", hex.EncodeToString(v.Hash[:]))
		fmt.Printf("bytes  := %d \n", v.Bytes)
		fmt.Println("")
	}
}

func HeaderPrint(v Slice) {
	fmt.Printf("hash   := %s \n", hex.EncodeToString(v.Hash[:]))
	//fmt.Printf("bytes  := %d \n", v.Bytes)
	fmt.Printf("len  := %d \n", v.Len)
	fmt.Println("")
}

type SliceOld struct {
	Nth    int
	Offset int64
	Hash   []byte
	Bytes  int
}

type Slice struct {
	Bytes []byte
	Hash  [32]byte
	Len   int
}

func readDataOld(fileName string) []SliceOld {

	data := make([]byte, 1<<20) //1Meg
	var slar []SliceOld

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; true; i++ {

		//get the file current offset
		//must do before read since read disturbs ofset
		off, err := file.Seek(0, os.SEEK_CUR)
		if err != nil {
			log.Fatal(err)
		}

		//will return err=EOF.. at EOF
		read, err := file.Read(data)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		//calculate the hash
		hash := sha256.Sum256(data[:read])

		slar = append(slar, SliceOld{Nth: i, Offset: off, Hash: hash[:], Bytes: read})
	}
	return slar
}

func readData(fileName string, chunk int) Slice {

	data := make([]byte, 1<<20) //1Meg

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	//get the file current offset
	//must do before read since read disturbs ofset
	_, err = file.Seek((int64((1 << 20) * chunk)), os.SEEK_CUR)
	if err != nil {
		log.Fatal(err)
	}

	//will return err=EOF.. at EOF
	ln, err := file.Read(data)
	if err == io.EOF {
		fmt.Println("EOF reached")
		panic(0)
	}

	if err != nil {
		log.Fatal(err)
	}

	//calculate the hash
	hash := sha256.Sum256(data[:ln])

	return Slice{data[0:ln], hash, ln}
}

func deleteFile(fileName string) error {
	return os.Remove(fileName)
}

func writeData(fileName string, chunk int, data2write []byte) bool {

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)

	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	//must do before read since read disturbs ofset
	_, err = file.Seek((int64((1 << 20) * chunk)), os.SEEK_CUR)
	if err != nil {
		log.Fatal(err)
		return false
	}

	_, err = file.Write(data2write)
	if err != nil {
		panic(0)
	}

	if err != nil {
		log.Fatal(err)
		return false
	}

	err = file.Sync()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//func main() {
//	HeaderPrintOld(readDataOld("testFile"))
//	HeaderPrint(readData("testFile", 1))
//}

//driver for writeData
//func main() {
//	s0 := readData("testFile", 1)
//	writeData("outFiel", 1, s0.Bytes)
//}
