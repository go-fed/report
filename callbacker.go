package report

import (
	"context"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams"
)

var _ pub.Callbacker = &NothingCallbacker{}

// Does nothing. In a real implementation, would handle a lot of details
// specific to that implementation.
type NothingCallbacker struct{}

func (n *NothingCallbacker) Create(c context.Context, s *streams.Create) error { return nil }
func (n *NothingCallbacker) Update(c context.Context, s *streams.Update) error { return nil }
func (n *NothingCallbacker) Delete(c context.Context, s *streams.Delete) error { return nil }
func (n *NothingCallbacker) Add(c context.Context, s *streams.Add) error       { return nil }
func (n *NothingCallbacker) Remove(c context.Context, s *streams.Remove) error { return nil }
func (n *NothingCallbacker) Like(c context.Context, s *streams.Like) error     { return nil }
func (n *NothingCallbacker) Block(c context.Context, s *streams.Block) error   { return nil }
func (n *NothingCallbacker) Follow(c context.Context, s *streams.Follow) error { return nil }
func (n *NothingCallbacker) Undo(c context.Context, s *streams.Undo) error     { return nil }
func (n *NothingCallbacker) Accept(c context.Context, s *streams.Accept) error { return nil }
func (n *NothingCallbacker) Reject(c context.Context, s *streams.Reject) error { return nil }
