package stripe

import (
	"github.com/pkg/errors"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/event"
)

// GetEvent retrieves the details of an event.
func GetEvent(eventID string) (*stripe.Event, error) {
	event, err := event.Get(eventID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "stripe: event")
	}

	return event, nil
}

// ListEvents returns a list of events going back up to 30 days.
func ListEvents() []*stripe.Event {
	var list []*stripe.Event

	i := event.List(nil)

	for i.Next() {
		e := i.Event()
		list = append(list, e)
	}

	return list
}
