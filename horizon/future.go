/*
an event horizon is a boundary beyond which events cannot affect an observer.

as it is with time, the order of events is not always linear.
*/
package horizon

import "log"

// The idea is simple....let a future event signal a past function to cause an effect in the present.
type Future struct {
	eventChan chan Event
	quitChan  chan struct{}

	onComFunc func(interface{})
	events    []Event

	onErrFunc func(interface{})

	onFinFunc func()

	signalCount int // could be useful

	// we have horizon.Einstein and horizon.Hawking
	// the mode defines whether or not the future should panic
	Mode Mode
}

// NewFuture creates a new future
func NewFuture(mode ...Mode) *Future {
	m := Einstein
	if len(mode) != 0 {
		m = mode[0]
	}
	f := &Future{eventChan: make(chan Event), events: make([]Event, 0), quitChan: make(chan struct{}), Mode: m}
	f.Signal()
	return f
}

// Events returns all the events returned from the future
func (f *Future) Events() []Event {
	return f.events
}

// Set passes the value received from the event that occured in the future
// into the event channel
func (f *Future) Set(e Event) {
	f.eventChan <- e
}

// RegisterComplete registers a function to be called when the future is complete
func (f *Future) RegisterComplete(futureFunc func(interface{})) {
	f.onComFunc = futureFunc
}

// RegisterError registers a function to be called when the future encounters an error
func (f *Future) RegisterError(futureFunc func(interface{})) {
	f.onErrFunc = futureFunc
}

// RegisterFinally registers a function to be called when the future is complete or encounters an error
func (f *Future) RegisterFinally(futureFunc func()) {
	f.onFinFunc = futureFunc
}

// Signal opens the horizon and allows the future to send events to the present
func (f *Future) Signal() {
	go func() {
	Loop:
		for {
			select {
			case e := <-f.eventChan:
				switch e.Type {
				case Complete:
					f.events = append(f.events, e)
				case Error:
					f.events = append(f.events, e)
				}
				f.signalCount++
			case <-f.quitChan:
				break Loop
			}
		}
	}()
}

// SignalComplete sends a signal to the future
// that the event has completed
func (f *Future) SignalComplete(value interface{}) {
	if f.isNil() {
		if f.Mode == Hawking {
			panic("future is destroyed")
		}
		log.Println("future is destroyed")
		return
	}

	if f.canSignal(Complete) {

		go func() {
			defer recovery(f)

			f.onComFunc(value)
			// handle error here -- only if user register a function for a future error event
			f.Set(Event{Type: Complete, Data: value})
		}()
	} else {
		if f.Mode == Hawking {
			panic("no function registered for future event [SignalComplete]")
		}
		log.Println("no function registered for future event [SignalComplete]")
		return
	}
}

// SignalError sends a signal to the future
// that the event has encountered an error
func (f *Future) SignalError(value interface{}) {
	if f.isNil() {
		if f.Mode == Hawking {
			panic("future is destroyed")
		}
		log.Println("future is destroyed")
		return
	}

	if f.canSignal(Error) {
		go func() {
			// any error inside this future handler, while
			// defer recovery(f) is active, yields an infinite panic loop
			// defer recovery(f)

			f.onErrFunc(value)
			// handle error here -- only if user register a function for a future error event
			f.Set(Event{Type: Error, Data: value})

			// signal finally
			f.signalFinally()
		}()
	} else {
		if f.Mode == Hawking {
			panic("no function registered for future event [SignalError]")
		}
		log.Println("no function registered for future event [SignalError]")
		return
	}
}

// SignalFinally sends a signal to the future
// that the event has completed or encountered an error
func (f *Future) signalFinally() {
	if f.isNil() {
		if f.Mode == Hawking {
			panic("future is destroyed")
		}
		log.Println("future is destroyed")
		return
	}

	if f.canSignal(Finally) {
		go func() {
			// any error that occurs inside a Finally handler
			// is out of bounds and won't be caught in the horizon
			// defer recovery(f)

			f.onFinFunc()
		}()
	} else {
		if f.Mode == Hawking {
			panic("no function registered for future event [SignalFinally]")
		}
		log.Println("no function registered for future event [SignalFinally]")
		return
	}
}

/*
SigalCount returns the number of signals sent to the future this is useful for debugging

note: it only accounts for the current timeline and has no knowledge of branched timelines created by f.Alter() **/
func (f *Future) SigalCount() int {
	return f.signalCount
}

/*
Alter changes the future, branching off to a new timeline

all signals are still registered but event histories are lost*/
func (f *Future) Alter() {
	f.events = nil
	f.signalCount = 0
}

// BlackHole closes the future, all signals are lost
// all branched timelines are pruned
// and the future is no longer usable
func (f *Future) BlackHole() {
	// prevent the event horizon from receiving any more signals
	close(f.quitChan)
	// keep the mode to allow for error handling
	// when accessing a black holed future
	*f = Future{Mode: f.Mode}
}

// isNil checks if the future is nil
// nil refers to a Future that has no values
// and is not usable
func (f *Future) isNil() bool {
	if f.eventChan == nil && f.events == nil && f.onComFunc == nil && f.quitChan == nil {
		return true
	}
	return false
}

// canSignal checks if the future has event handlers registered
// and is not nil
func (f *Future) canSignal(t EventType) bool {
	if f.isNil() {
		return false
	}

	switch t {
	case Complete:
		return f.onComFunc != nil
	case Error:
		return f.onErrFunc != nil
	case Finally:
		return f.onFinFunc != nil
	}

	return false
}
