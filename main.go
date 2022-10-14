package main

import (
	"fmt"
	"time"

	"github.com/opensaucerer/wormhole/spacetime"
)

func main() {

	f := spacetime.NewFuture()
	f.RegisterComplete(func(data interface{}) {
		fmt.Println("Received data from future --> ", data.(string))
	})

	go func() {
		<-time.After(4 * time.Second)
		f.SignalComplete("completed after 4 seconds")
	}()

	go func() {
		<-time.After(2 * time.Second)
		f.SignalComplete("completed after 2 seconds")
	}()

	f.SignalComplete("another hello world")

	time.Sleep(5 * time.Second)
}
