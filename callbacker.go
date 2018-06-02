package report

import (
	"context"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams"
)

var _ pub.Callbacker = &nothingCallbacker{}

// Does nothing. In a real implementation, would handle a lot of details
// specific to that implementation.
type nothingCallbacker struct{}

func (n *nothingCallbacker) Create(c context.Context, s *streams.Create) error { return nil }
func (n *nothingCallbacker) Update(c context.Context, s *streams.Update) error { return nil }
func (n *nothingCallbacker) Delete(c context.Context, s *streams.Delete) error { return nil }
func (n *nothingCallbacker) Add(c context.Context, s *streams.Add) error       { return nil }
func (n *nothingCallbacker) Remove(c context.Context, s *streams.Remove) error { return nil }
func (n *nothingCallbacker) Like(c context.Context, s *streams.Like) error     { return nil }
func (n *nothingCallbacker) Block(c context.Context, s *streams.Block) error   { return nil }
func (n *nothingCallbacker) Follow(c context.Context, s *streams.Follow) error { return nil }
func (n *nothingCallbacker) Undo(c context.Context, s *streams.Undo) error     { return nil }
func (n *nothingCallbacker) Accept(c context.Context, s *streams.Accept) error { return nil }
func (n *nothingCallbacker) Reject(c context.Context, s *streams.Reject) error { return nil }
