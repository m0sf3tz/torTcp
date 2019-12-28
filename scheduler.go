// In this example we'll look at how to implement
// a _worker pool_ using goroutines and channels.

package main

import (
	//"fmt"
	"math/rand"
	"time"
)

func randX(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

// concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding
// results on `results`.
func worker(id int, jobs <-chan chunkStructWithIp, results chan<- chunkStruct) {
	for j := range jobs {
		r := getChunk(j.fileName, j.chunk, j.ip)

		if r == true {
			j.done = true
		} else {
			j.done = false
		}
		results <- chunkStruct{j.chunk, j.done, j.fileName}
	}
}

type chunkStruct struct {
	chunk    int
	done     bool
	fileName string
}

type chunkStructWithIp struct {
	chunk    int
	done     bool
	fileName string
	ip       string
}

//enqueer
func dispatcher(jobs chan<- chunkStructWithIp, jb []chunkStruct) {

	//enque all the work we have then close the
	//channel to let the works know there is no work
	//left
	for _, j := range jb {

		var t chunkStructWithIp

		t.chunk = j.chunk
		t.done = j.done
		t.fileName = j.fileName
		t.ip = al.seeders[randX(len(al.seeders))].SeedIp
		newIp := RemovePort(t.ip)

		t.ip = newIp + ":8082"
		jobs <- t
	}
	close(jobs)
}

func scheduler(jb []chunkStruct, chunksLeft int, seeders int) []chunkStruct {
	workLeft := make([]chunkStruct, 0, 100)

	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan chunkStructWithIp, seeders)
	results := make(chan chunkStruct, seeders)

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 0; w < 2; w++ {
		go worker(w, jobs, results)
	}

	go dispatcher(jobs, jb)

	// Finally we collect all the results of the work.
	// This also ensures that the worker goroutines have
	// finished.
	for a := 0; a < chunksLeft; a++ {
		x := <-results
		if x.done == false {
			workLeft = append(workLeft, x)
		}
	}
	//fmt.Println(workLeft)
	return workLeft
}

func fetchTop(jb []chunkStruct, seeders int) {

	for {
		jb_ret := (scheduler(jb, len(jb), seeders))

		jb = jb_ret

		if len(jb) == 0 {
			break
		}
	}
}
