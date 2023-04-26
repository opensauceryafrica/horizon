package horizon

type Radius struct {
	X, Y, Z float64
}

// the radius at which an horizon blocks forever
var SchwarzschildRadius Radius = Radius{0, 0, 0}

// the radius at which a blocked horizon becomes unblocked
var EinsteinRadius Radius = Radius{1, 1, 1}

// Schwarzschild creates a warps in the horizon such that the runtime blocks
// until signal off. This signal off simply changes the coordinate away from
// the schwarzschild radius.
func Schwarzschild(future *Future) {
	if !future.isNil() {
		future.radius = SchwarzschildRadius
	}
	for {
		if future.radius != SchwarzschildRadius {
			break
		}
	}
}

// Openheimer shifts the horion radius away from the schwarzschild radius
// and towards the einstein radius. This allows the future to be unblocked
// and the runtime to continue.
func Openheimer(future *Future) {
	if !future.isNil() {
		future.radius = EinsteinRadius
	}
}
