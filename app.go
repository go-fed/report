package report

import (
	"context"
	"crypto"
	"fmt"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/vocab"
	"github.com/go-fed/httpsig"
	"net/http"
	"net/url"
	"sync"
)

var _ pub.Application = &App{}
var _ pub.SocialAPI = &App{}
var _ pub.FederateAPI = &App{}
var _ pub.SocialApplication = &App{}
var _ pub.FederateApplication = &App{}
var _ pub.SocialFederateApplication = &App{}

type lockObj struct {
	obj pub.PubObject
	mu  *sync.RWMutex
}

// App shows the basic mechanics for a single-user, non-permanent, dummy server.
type App struct {
	scheme    string
	host      string
	newPath   string
	db        map[string]lockObj
	dbMu      *sync.RWMutex
	actor     *vocab.Object
	actorMu   *sync.RWMutex
	actorURL  *url.URL
	inboxURL  *url.URL
	outboxURL *url.URL
	id        int
	idMu      *sync.Mutex
	pubKey    crypto.PublicKey
	privKey   crypto.PrivateKey
}

func NewApp(scheme, host, newPath string, actorURL, inboxURL, outboxURL *url.URL, pubKey crypto.PublicKey, privKey crypto.PrivateKey, actor *vocab.Object) *App {
	return &App{
		scheme:    scheme,
		host:      host,
		newPath:   newPath,
		db:        make(map[string]lockObj),
		dbMu:      &sync.RWMutex{},
		actor:     actor,
		actorMu:   &sync.RWMutex{},
		actorURL:  actorURL,
		inboxURL:  inboxURL,
		outboxURL: outboxURL,
		id:        1,
		idMu:      &sync.Mutex{},
		pubKey:    pubKey,
		privKey:   privKey,
	}
}

func (a *App) Owns(c context.Context, id *url.URL) bool {
	return id.Host == a.host
}

func (a *App) Get(c context.Context, id *url.URL, rw pub.RWType) (pub.PubObject, error) {
	has, err := a.Has(c, id)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("%s not found", id)
	}
	a.dbMu.RLock()
	p, ok := a.db[id.String()]
	a.dbMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%s not found", id)
	}
	switch rw {
	case pub.Read:
		p.mu.RLock()
		go func() {
			<-c.Done()
			p.mu.RUnlock()
		}()
	case pub.ReadWrite:
		p.mu.Lock()
		go func() {
			<-c.Done()
			p.mu.Unlock()
		}()
	default:
		return nil, fmt.Errorf("unrecognized pub.RWType: %v", rw)
	}
	return p.obj, nil
}

func (a *App) GetAsVerifiedUser(c context.Context, id, authdUser *url.URL, rw pub.RWType) (pub.PubObject, error) {
	return a.Get(c, id, rw)
}

func (a *App) Has(c context.Context, id *url.URL) (bool, error) {
	_, ok := a.db[id.String()]
	return ok, nil
}

func (a *App) Set(c context.Context, o pub.PubObject) error {
	if id := o.GetId(); id == nil {
		return fmt.Errorf("id is nil")
	} else if *id == *a.outboxURL {
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetOutboxOrderedCollection(oc)
		return nil
	} else if *id == *a.inboxURL {
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetInboxOrderedCollection(oc)
		return nil
	} else {
		a.dbMu.RLock()
		a.db[id.String()] = lockObj{
			obj: o,
			mu:  &sync.RWMutex{},
		}
		a.dbMu.RUnlock()
		return nil
	}
}

func (a *App) GetInbox(c context.Context, r *http.Request, rw pub.RWType) (vocab.OrderedCollectionType, error) {
	if *r.URL == *a.inboxURL {
		a.actorMu.RLock()
		defer a.actorMu.RUnlock()
		return a.actor.GetInboxOrderedCollection(), nil
	}
	return nil, fmt.Errorf("no inbox for url %s", r.URL)
}

func (a *App) GetOutbox(c context.Context, r *http.Request, rw pub.RWType) (vocab.OrderedCollectionType, error) {
	if *r.URL == *a.outboxURL {
		a.actorMu.RLock()
		defer a.actorMu.RUnlock()
		return a.actor.GetOutboxOrderedCollection(), nil
	}
	return nil, fmt.Errorf("no outbox for url %s", r.URL)
}

func (a *App) NewId(c context.Context, t pub.Typer) *url.URL {
	a.idMu.Lock()
	id := a.id
	a.id++
	a.idMu.Unlock()
	withoutTrailingSlash := a.newPath
	if a.newPath[len(a.newPath)-1] == '/' {
		withoutTrailingSlash = a.newPath[:len(a.newPath)-1]
	}
	return &url.URL{
		Scheme: a.scheme,
		Host:   a.host,
		Path:   fmt.Sprintf("%s/%d", withoutTrailingSlash, id),
	}
}

func (a *App) GetPublicKey(c context.Context, publicKeyId string) (pubKey crypto.PublicKey, algo httpsig.Algorithm, user *url.URL, err error) {
	return nil, httpsig.RSA_SHA256, nil, fmt.Errorf("not implemented: GetPublicKey")
}

func (a *App) CanAdd(c context.Context, o vocab.ObjectType, t vocab.ObjectType) bool {
	return true
}

func (a *App) CanRemove(c context.Context, o vocab.ObjectType, t vocab.ObjectType) bool {
	return true
}

func (a *App) ActorIRI(c context.Context, r *http.Request) (*url.URL, error) {
	if *r.URL == *a.inboxURL || *r.URL == *a.outboxURL {
		return a.actorURL, nil
	}
	return nil, fmt.Errorf("no actor for url %s", r.URL)
}

func (a *App) GetSocialAPIVerifier(c context.Context) pub.SocialAPIVerifier {
	// TODO: OAuth 2
	return nil
}

func (a *App) GetPublicKeyForOutbox(c context.Context, publicKeyId string, boxIRI *url.URL) (crypto.PublicKey, httpsig.Algorithm, error) {
	if boxIRI != a.outboxURL {
		return nil, httpsig.RSA_SHA256, fmt.Errorf("unknown outbox url %s", boxIRI)
	} else if publicKeyId != a.actorURL.String() {
		return nil, httpsig.RSA_SHA256, fmt.Errorf("unknown public key id %q", publicKeyId)
	}
	return a.pubKey, httpsig.RSA_SHA256, nil
}

func (a *App) OnFollow(c context.Context, s *streams.Follow) pub.FollowResponse {
	return pub.AutomaticAccept
}

func (a *App) Unblocked(c context.Context, actorIRIs []*url.URL) error {
	return nil
}

func (a *App) FilterForwarding(c context.Context, activity vocab.ActivityType, iris []*url.URL) ([]*url.URL, error) {
	// Do NOT do this in real implementations. This turns the server into a
	// spambot. See the documentation in go-fed/activity/pub.
	return iris, nil
}

func (a *App) NewSigner() (httpsig.Signer, error) {
	s, _, err := httpsig.NewSigner([]httpsig.Algorithm{httpsig.RSA_SHA256}, nil, httpsig.Signature)
	return s, err
}

func (a *App) PrivateKey(boxIRI *url.URL) (privKey crypto.PrivateKey, pubKeyId string, err error) {
	return a.privKey, a.actorURL.String(), nil
}
