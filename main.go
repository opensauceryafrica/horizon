package main

import (
	"fmt"
	"log"
	"time"

	"github.com/opensaucerer/wormhole/horizon"
)

func main() {

	f := horizon.NewFuture(horizon.Einstein)

	f.RegisterComplete(func(data interface{}) {
		log.Println("Received data from future --> ", data.(string))

		// an error in this future handler with send a trigger to the future error handler
		// panic("We should not be here")

		if data.(string) == "completed after 4 seconds" {
			horizon.Openheimer(f)
		}
	})

	// f.RegisterError(func(data interface{}) {
	// 	log.Println("Received error from future --> ", data.(string))

	// })

	f.RegisterFinally(func() {
		log.Println("Future is complete")
	})

	go func() {
		<-time.After(4 * time.Second)
		f.SignalComplete("completed after 4 seconds")
	}()

	go func() {
		<-time.After(2 * time.Second)
		f.SignalComplete("completed after 2 seconds")
	}()

	for i := 0; i < 10; i++ {

		f.SignalComplete(" hello world: " + fmt.Sprintf("%d", i))
	}

	horizon.Schwarzschild(f)

	// <-time.After(10 * time.Second)
	// log.Println(f.Events(), f.SigalCount())

	// f.Alter()

	// log.Println(f.Events())

	// f.SignalComplete("hello world")

	// time.Sleep(5 * time.Second)

	// log.Println(f.Events())

	// f.BlackHole()

	// time.Sleep(5 * time.Second)

	// f.SignalComplete("future destroyed")

	// time.Sleep(5 * time.Second)
}
