package main

import (
	"fmt"
	"time"
)

/*
A pipeline is simply functions that process data sequntially
Given example of simple conccurrent pipeline
genrator :  a function that coverts data to steams
stages : a function that takes bussines logic converts it to a streams
*/
func generator(done <-chan interface{}) <-chan int {
	stream := make(chan int)

	go func() {
		defer close(stream)
		data := 10
		for {
			select {
			case <-done:
				return
			case stream <- data:
				data += 10
			}

		}
	}()
	return stream
}

func Stage(done <-chan interface{}, dataStream <-chan int, logic func(data int) int, name string) <-chan int {
	outpuStream := make(chan int)
	go func() {
		defer close(outpuStream)
		for data := range dataStream {
			select {
			case <-done:
				return
			case outpuStream <- logic(data):
				fmt.Printf("\nStage %s data: %d", name, logic(data))
				time.Sleep(time.Second * 2)
			}

		}
	}()
	return outpuStream
}

func main() {
	done := make(chan interface{})

	// create a stream
	stream := generator(done)

	add := func(data int) int {
		return data + 10
	}

	mul := func(data int) int {
		return data * 10
	}

	Stage1 := Stage(done, stream, add, "Add Stage")

	Stage2 := Stage(done, Stage1, mul, "Mul Stage")
	count := 0

	go func() {
		time.Sleep(time.Second * 10)
		close(done)
	}()
	for data := range Stage2 {
		count++
		fmt.Printf("\n ::::  %d Result %d", count, data)
	}

}
