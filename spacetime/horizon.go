package spacetime

type Future struct {
	completeChan  chan interface{}
	onComFunc     interface{}
	completeEvent []interface{}
	signalCount   int // could be useful
}

func NewFuture() *Future {
	return &Future{completeChan: make(chan interface{})}
}

func (f *Future) GetCompleteEventFromFuture(signalId int) interface{} {
	if signalId < f.signalCount {
		return f.completeEvent[signalId]
	}
	return nil
}

func (f *Future) GetCompleteEventsFromFuture() []interface{} {
	return f.completeEvent
}

func (f *Future) Set(value interface{}, future string) {
	switch future {
	case "complete":
		f.completeChan <- value
	default:
	}
}

func (f *Future) RegisterComplete(futureFunc interface{}) {
	f.onComFunc = futureFunc
}

func (f *Future) Signal() {
	// maybe this should be a blocking call?
	go func() {
	Loop:
		for {
			select {
			case e := <-f.completeChan:
				f.completeEvent = append(f.completeEvent, e)
				f.signalCount++
				break Loop
			default:
				break Loop
			}
		}
	}()
}

func (f *Future) SignalComplete(value interface{}) {
	if f.onComFunc != nil {
		go func() {
			f.onComFunc.(func(interface{}))(value)
			// handle error here -- only if user register a function for a future error event
			f.Set(value, "complete")
		}()
		f.Signal()
	} else {
		panic("no function registered for future event [SignalComplete]")
	}
}

func (f *Future) SigalCount() int {
	return f.signalCount
}
