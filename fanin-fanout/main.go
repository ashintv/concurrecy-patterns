package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

/*
	FAN IN - FAN OUT

	A pattern which can be used when order of execution of a stream deosent matter
	DoWork()  any work

	fan-out takes the data and spin ups go routines for doing the work
	fan-in combines goroutines (workers) for result to a single channels
*/

// so work is a function that simualetes a realworld work
func doWork(done <-chan interface{}) <-chan interface{} {

	workStream := make(chan interface{})

	go func() {
		num := rand.Intn(5)
		log.Printf("Routine Write after %ds" , num)
		defer close(workStream)
		select {
		case <-done:
		case <-time.After(time.Duration(num) * time.Second):

			workStream <- num
		}

	}()
	// adding a sleep for simulating work
	return workStream
}

func fanout(done <-chan interface{}) []<-chan interface{} {
	log.Printf("Spinning 3 ")
	numWorkers := 3
	workers := make([]<-chan interface{}, numWorkers)
	for i := range numWorkers {
		workers[i] = doWork(done)
	}
	return workers
}

func fanin(done <-chan interface{}, channels []<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multPlexedStream := make(chan interface{})
	log.Printf("Combing data")
	multiplexer := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multPlexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))

	for _, c := range channels {
		go multiplexer(c)
	}

	go func() {
		wg.Wait()
		close(multPlexedStream)
	}()
	return multPlexedStream
}

func main() {

	cntlr := make(chan interface{})
	defer close(cntlr)
	var done <-chan interface{}
	done = cntlr
	resultChanneles := fanout(done)
	resultStream := fanin(done, resultChanneles)
	for data := range resultStream {
		log.Printf("Data Recieved %d", data)
	}

}
