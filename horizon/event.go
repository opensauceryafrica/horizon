package horizon

// Event represents a piece of information from the future
// that is sent to the present via the past.
type Event struct {
	Type EventType
	Data interface{}
}

// EventType is a type that represents the type of event
type EventType string

const (
	// Complete is the event type for a complete event
	Complete EventType = "complete"

	// Error is the event type for an error event
	Error EventType = "error"

	// Finally is the event type for a finally event
	Finally EventType = "finally"
)
