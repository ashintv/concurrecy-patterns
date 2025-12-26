package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*

	In many situations we may have to do or operation in a no.of channels
	for that creating a single or channel that can take n no.of channels a return a single channel
	the single channels get called (done) whenever any of the input channel get read

*/

func or(channels ...<-chan interface{}) <-chan interface{} {

	switch len(channels) {
	case 0:
		return nil

	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})

	defer close(orDone)
	go func() {
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}

		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}

func main() {
	sig := func() <-chan interface{} {
		c := make(chan interface{})

		go func() {
			defer close(c)
			sleep := rand.Intn(5) + 1
			fmt.Println("Sleeping for", sleep, "seconds")
			time.Sleep(time.Duration(sleep) * time.Second)

		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(),
		sig(),
		sig(),
		sig(),
	)
	fmt.Println("Done after:", time.Since(start))
}
