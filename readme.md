# horizon

Horizon creates a boundary beyond which events can no longer affect the observer, however, as it is with time, the order of events isn't always linear.

The idea is simple...let a future event signal a past function to cause an effect in the present.

## Why Horizon?

Goroutines are arguably a transgenerational piece of programming art, but they can be quite difficult to manage. Of course, you have channels but it's impossible to deny that you haven't caused the runtime to block by reading from an unbuffered channel or even a situation that's worse, a deadlock.

Rather than having to deal with the complexity of channels, Horizon provides a simple interface to manage goroutines and their events in a way that makes it much flexible for the most common use cases. It's a simple implementation of the [Future](https://en.wikipedia.org/wiki/Futures_and_promises) pattern.

## Installation

```shell
go get github.com/opensaucerer/horizon
```

## Usage

### Sending emails to users with retries and retry timeout

```go
package main

import (
	"log"
	"time"

	"github.com/opensaucerer/wormhole/horizon"
)

func main() {

	var sent, failed int
	retrial := make(map[string]int)

	future := horizon.NewFuture(horizon.Einstein)

	future.RegisterComplete(func(data interface{}) {
		mail := data.(service.Email)

		log.Printf("EmailNotification: notification for %s successfully sent to %s", mail.MessageKey, mail.To)

		sent++

		if sent+failed == len(recipients) {
			log.Printf("EmailNotification: %d notifications sent, %d failed", sent, failed)

			// unblock the goroutine that called this function
			horizon.Openheimer(future)
		}
	})

	future.RegisterError(func(data interface{}) {
		mail := data.(service.Email)

		log.Printf("EmailNotification: notification for %s failed to send to %s --- Retrying", mail.MessageKey, mail.To)

		if retrial[mail.To] > 2 {
			failed++
		} else {
			retrial[mail.To]++
			go func(m service.Email) {
				time.Sleep(5 * time.Second)
				err := service.SendEmailForNotification(m.To, m.MessageKey, m.Replacements)
				if err != nil {
					future.SignalError(m)
				} else {
					future.SignalComplete(m)
				}
			}(mail)
		}

		if sent+failed == len(recipients) {

			log.Printf("EmailNotification: %d notifications sent, %d failed", sent, failed)

			// unblock the main routine
			horizon.Openheimer(future)
		}
	})

	for _, mail := range recipients {
		go func(m service.Email) {
			time.Sleep(5 * time.Second)
			err := service.SendEmailForNotification(m.To, m.MessageKey, m.Replacements)
			if err != nil {
				future.SignalError(m)
			} else {
				future.SignalComplete(m)
			}
		}(mail)
	}

	// block the main routine
	horizon.Schwarzschild(future)
}
```

## Functions

### NewFuture

Creates a new future along the horizon. It optionally takes a `mode` argument which can be one of the following:

- `horizon.Einstein` - This ensures that the future will not panic in any case of timeline interference.

- `horizon.Hawking` - This allows the future to panic in the case of timeline interference.

### RegisterComplete

Takes a function of the form `func(interface{}){}` as an argument and registers it as a callback for when the future is signaled to as complete.

### RegisterError

Takes a function of the form `func(interface{}){}` as an argument and registers it as a callback for when the future is signaled to as an error.

### RegisterFinally

Gets called when the future is signaled to as complete or error. It takes a function of the form `func(){}` as an argument. It's clear that no data is sent to the past via the horizon when this this is called.

### SignalComplete

Takes a single argument of type `interface{}` and signals the future as complete. The argument is sent to the past via the horizon.

### SignalError

Takes a single argument of type `interface{}` and signals the future as error. The argument is sent to the past via the horizon.

### Schwarzschild

Schwarzschild creates a warps in the horizon such that the runtime blocks until signaled off. This signal off simply changes the coordinate away from the schwarzschild radius.

### Openheimer

Openheimer shifts the horion radius away from the schwarzschild radius and towards the einstein radius. This allows the future's horizon to be unblocked and the runtime to continue.

### SignalCount

Returns the number of signals sent to the future. This could be useful for debugging

### Alter

Changes the future, branching off to a new timeline. All signals are still registered but event histories are lost

### BlackHole

Closes the future. All signals are lost, all branched timelines are pruned and the future is no longer usable

### Events

Returns all the events returned from the future via the horizon
