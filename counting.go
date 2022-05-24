package ingnovus

import "fmt"

type counting struct {
	*event
}
type Counting interface {
	Event
}

func ParseCounting(e Event) (Counting, error) {
	if e.Type() == CountingEvent {
		v, ok := e.(Counting)
		if !ok {
			return nil, fmt.Errorf("event is not a CountingEvent, type: %v", e.Type())
		}
		return v, nil
	}
	return nil, fmt.Errorf("event is not a CountingEvent, type: %v", e.Type())
}
