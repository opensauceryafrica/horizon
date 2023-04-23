# wormhole

The idea is simple....let a future event signal a past function to cause an effect in the present.

```go
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

	<-time.After(10 * time.Second)
	log.Println(f.Events(), f.SigalCount())

	f.Alter()

	log.Println(f.Events())

	f.SignalComplete("hello world")

	time.Sleep(5 * time.Second)

	log.Println(f.Events())

	f.BlackHole()

	time.Sleep(5 * time.Second)

	f.SignalComplete("future destroyed")

	time.Sleep(5 * time.Second)
}
```

#### result

```shell
Received data from future -->  another hello world
Received data from future -->  completed after 2 seconds
Received data from future -->  completed after 4 seconds
```

## TODO

- [x] basic Futures implementation
- [x] handle error in Futures
- [x] `signalFinally` in Futures (called immediately after SignalComplete or SignalError)
- [x] `SingalError` in Futures

> come in and let's hack!
