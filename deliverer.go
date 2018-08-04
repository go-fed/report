package report

import (
	"github.com/go-fed/activity/pub"
	"log"
	"net/url"
)

var _ pub.Deliverer = &syncDeliverer{}

// Synchronously tries to deliver, immediately. Real applications may want to
// rate limit, back off, and retry across downtimes.
type syncDeliverer struct{}

func (s *syncDeliverer) Do(b []byte, to *url.URL, toDo func(b []byte, u *url.URL) error) {
	log.Printf("Delivering to: %s", to)
	err := toDo(b, to)
	if err != nil {
		log.Print(err)
	}
}
