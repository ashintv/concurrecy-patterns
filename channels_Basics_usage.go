package main

import "fmt"

/*
	How to use the channels properly?
	Implementatiing proper (Safer way of channels )
	Based on following rules
		1.  Every channel has a owner (go rotines that creats channel) and owners ship can be transferd
		2.	The owners is responsible for closing the channels
		3.	The data creation or write to Channel shouls close as possibe to the creation of channels
*/

func main() {
	// producer functions is the function that creates the channels
	// since producer is the
	producer := func() <-chan int {

		//create the channels
		var DataStream chan int
		DataStream = make(chan int, 10)

		// create a data stream (Simulation)
		go func() {

			//rule 2
			defer close(DataStream)
			for i := range 10 {
				DataStream <- i
			}
		}()
		return DataStream
	}

	// Consumer

	consumer := func(Stream <-chan int) {
		for data := range Stream {
			fmt.Printf("\nData Read %d", data)
		}
		fmt.Println("\nStream ended")
	}

	Stream := producer()
	consumer(Stream)
}
