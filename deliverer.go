package report

import (
	"github.com/go-fed/activity/pub"
	"log"
	"net/url"
)

var _ pub.Deliverer = &SyncDeliverer{}

// Synchronously tries to deliver, immediately. Real applications may want to
// rate limit, back off, and retry across downtimes.
type SyncDeliverer struct{}

func (s *SyncDeliverer) Do(b []byte, to *url.URL, toDo func(b []byte, u *url.URL) error) {
	err := toDo(b, to)
	if err != nil {
		log.Print(err)
	}
}
