# wormhole

The idea is simple....let a future event signal a past function to cause an effect in the present.

```go
package main

import (
	"fmt"
	"time"

	"github.com/opensaucerer/wormhole/horizon"
)

func main() {

	f := horizon.NewFuture()
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
