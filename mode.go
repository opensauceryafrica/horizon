package horizon

const (
	// Einstein is the default mode for the future. This ensures that the future will not panic in any case of timeline interference.
	Einstein Mode = iota

	// Hawking is the mode for the future that allows the future to panic in the case of timeline interference.
	Hawking
)

// Mode defines whether or not the future should panic
type Mode int

// Int returns the integer value of the mode
func (m Mode) Int() int {
	return int(m)
}
