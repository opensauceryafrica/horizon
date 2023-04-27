package horizon

// recovery traverses the path of the future and recovers from any panics that may have occured.
func recovery(f *Future) {
	if r := recover(); r != nil {
		if f.canSignal(Error) {
			// signal error
			f.SignalError(r)
		}
	} else {
		if f.canSignal(Finally) {
			f.signalFinally()
		}
	}
}
