package report

import (
	"context"
	"crypto"
	"fmt"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/vocab"
	"github.com/go-fed/httpsig"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var _ pub.Application = &app{}
var _ pub.SocialAPI = &app{}
var _ pub.FederateAPI = &app{}
var _ pub.SocialApplication = &app{}
var _ pub.FederateApplication = &app{}
var _ pub.SocialFederateApplication = &app{}

type lockKeyType string

type lockObj struct {
	obj pub.PubObject
	mu  *sync.RWMutex
	who int
}

// app shows the basic mechanics for a single-user, non-permanent, dummy server.
type app struct {
	scheme       string
	host         string
	newPath      string
	db           map[string]*lockObj
	dbMu         *sync.RWMutex
	actor        *vocab.Person
	actorURL     *url.URL
	inboxURL     *url.URL
	outboxURL    *url.URL
	followingURL *url.URL
	followersURL *url.URL
	likedURL     *url.URL
	inbox        vocab.OrderedCollectionType
	inboxMu      *sync.RWMutex
	outbox       vocab.OrderedCollectionType
	outboxMu     *sync.RWMutex
	following    vocab.OrderedCollectionType
	followingMu  *sync.RWMutex
	followers    vocab.OrderedCollectionType
	followersMu  *sync.RWMutex
	liked        vocab.OrderedCollectionType
	likedMu      *sync.RWMutex
	id           int
	idMu         *sync.Mutex
	pubKey       crypto.PublicKey
	privKey      crypto.PrivateKey
	verifier     pub.SocialAPIVerifier
}

func newApp(scheme, host, newPath string, actorURL, inboxURL, outboxURL, followingURL, followersURL, likedURL *url.URL, pubKey crypto.PublicKey, privKey crypto.PrivateKey, actor *vocab.Person, verifier pub.SocialAPIVerifier) *app {
	inbox := &vocab.OrderedCollection{}
	inbox.SetId(inboxURL)
	outbox := &vocab.OrderedCollection{}
	outbox.SetId(outboxURL)
	following := &vocab.OrderedCollection{}
	following.SetId(followingURL)
	followers := &vocab.OrderedCollection{}
	followers.SetId(followersURL)
	liked := &vocab.OrderedCollection{}
	liked.SetId(likedURL)
	return &app{
		scheme:       scheme,
		host:         host,
		newPath:      newPath,
		db:           make(map[string]*lockObj),
		dbMu:         &sync.RWMutex{},
		actor:        actor,
		actorURL:     actorURL,
		inboxURL:     inboxURL,
		outboxURL:    outboxURL,
		followingURL: followingURL,
		followersURL: followersURL,
		likedURL:     likedURL,
		inbox:        inbox,
		inboxMu:      &sync.RWMutex{},
		outbox:       outbox,
		outboxMu:     &sync.RWMutex{},
		following:    following,
		followingMu:  &sync.RWMutex{},
		followers:    followers,
		followersMu:  &sync.RWMutex{},
		liked:        liked,
		likedMu:      &sync.RWMutex{},
		id:           1,
		idMu:         &sync.Mutex{},
		pubKey:       pubKey,
		privKey:      privKey,
		verifier:     verifier,
	}
}

func (a *app) Owns(c context.Context, id *url.URL) bool {
	log.Printf("Owns: %s", id)
	return id.Host == a.host
}

func (a *app) Get(c context.Context, id *url.URL, rw pub.RWType) (pub.PubObject, error) {
	log.Printf("Getting: %s", id)
	has, err := a.Has(c, id)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("%s not found", id)
	}
	if *id == *a.actorURL {
		return a.actor, nil
	} else if *id == *a.inboxURL {
		a.inboxMu.RLock()
		defer a.inboxMu.RUnlock()
		return a.inbox, nil
	} else if *id == *a.outboxURL {
		a.outboxMu.RLock()
		defer a.outboxMu.RUnlock()
		return a.outbox, nil
	} else if *id == *a.followingURL {
		a.followingMu.RLock()
		defer a.followingMu.RUnlock()
		return a.following, nil
	} else if *id == *a.followersURL {
		a.followersMu.RLock()
		defer a.followersMu.RUnlock()
		return a.followers, nil
	} else if *id == *a.likedURL {
		a.likedMu.RLock()
		defer a.likedMu.RUnlock()
		return a.liked, nil
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
		defer p.mu.RUnlock()
	case pub.ReadWrite:
		who := c.Value(lockKeyType("lockKey")).(int)
		if p.who != who {
			log.Printf("locking %s", id)
			p.mu.Lock()
			p.who = who
			go func() {
				<-c.Done()
				if p.who == who {
					log.Printf("unlocking %s", id)
					p.mu.Unlock()
				}
			}()
		}
	default:
		return nil, fmt.Errorf("unrecognized pub.RWType: %v", rw)
	}
	return p.obj, nil
}

func (a *app) GetAsVerifiedUser(c context.Context, id, authdUser *url.URL, rw pub.RWType) (pub.PubObject, error) {
	log.Printf("GetAsVerifiedUser: %s", id)
	return a.Get(c, id, rw)
}

func (a *app) Has(c context.Context, id *url.URL) (bool, error) {
	log.Printf("Has: %s", id)
	if *id == *a.actorURL || *id == *a.inboxURL || *id == *a.outboxURL || *id == *a.followingURL || *id == *a.followersURL || *id == *a.likedURL {
		return true, nil
	}
	_, ok := a.db[id.String()]
	return ok, nil
}

func (a *app) Set(c context.Context, o pub.PubObject) error {
	b, _ := o.Serialize()
	log.Printf("Setting: %s", b)
	if id := o.GetId(); id == nil {
		return fmt.Errorf("id is nil")
	} else if *id == *a.outboxURL {
		a.outboxMu.Lock()
		defer a.outboxMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.outbox = oc
		return nil
	} else if *id == *a.inboxURL {
		a.inboxMu.Lock()
		defer a.inboxMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.inbox = oc
		return nil
	} else if *id == *a.followingURL {
		a.followingMu.Lock()
		defer a.followingMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.following = oc
		return nil
	} else if *id == *a.followersURL {
		a.followersMu.Lock()
		defer a.followersMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.followers = oc
		return nil
	} else if *id == *a.likedURL {
		a.likedMu.Lock()
		defer a.likedMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.liked = oc
		return nil
	} else {
		a.dbMu.Lock()
		if v, ok := a.db[id.String()]; ok {
			a.dbMu.Unlock()
			who := c.Value(lockKeyType("lockKey")).(int)
			vWho := v.who
			// TODO: Use sync.Cond
			if vWho == 0 || vWho != who {
				log.Printf("locking %s", id)
				v.mu.Lock()
				v.who = who
			}
			v.obj = o
			v.who = 0
			log.Printf("unlocking %s", id)
			v.mu.Unlock()
		} else {
			a.db[id.String()] = &lockObj{
				obj: o,
				mu:  &sync.RWMutex{},
			}
			a.dbMu.Unlock()
		}
		return nil
	}
}

func (a *app) GetInbox(c context.Context, r *http.Request, rw pub.RWType) (vocab.OrderedCollectionType, error) {
	log.Printf("GetInbox: %s", r.URL)
	if *r.URL == *a.inboxURL {
		a.inboxMu.RLock()
		defer a.inboxMu.RUnlock()
		return a.inbox, nil
	}
	return nil, fmt.Errorf("no inbox for url %s", r.URL)
}

func (a *app) GetOutbox(c context.Context, r *http.Request, rw pub.RWType) (vocab.OrderedCollectionType, error) {
	log.Printf("GetOutbox: %s", r.URL)
	if *r.URL == *a.outboxURL {
		a.outboxMu.RLock()
		defer a.outboxMu.RUnlock()
		return a.outbox, nil
	}
	return nil, fmt.Errorf("no outbox for url %s", r.URL)
}

func (a *app) NewId(c context.Context, t pub.Typer) *url.URL {
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

func (a *app) GetPublicKey(c context.Context, publicKeyId string) (pubKey crypto.PublicKey, algo httpsig.Algorithm, user *url.URL, err error) {
	return nil, httpsig.RSA_SHA256, nil, fmt.Errorf("not implemented: GetPublicKey")
}

func (a *app) CanAdd(c context.Context, o vocab.ObjectType, t vocab.ObjectType) bool {
	return true
}

func (a *app) CanRemove(c context.Context, o vocab.ObjectType, t vocab.ObjectType) bool {
	return true
}

func (a *app) ActorIRI(c context.Context, r *http.Request) (*url.URL, error) {
	log.Printf("ActorIRI: %s", r.URL)
	if *r.URL == *a.inboxURL || *r.URL == *a.outboxURL {
		return a.actorURL, nil
	}
	return nil, fmt.Errorf("no actor for url %s", r.URL)
}

func (a *app) GetSocialAPIVerifier(c context.Context) pub.SocialAPIVerifier {
	return a.verifier
}

func (a *app) GetPublicKeyForOutbox(c context.Context, publicKeyId string, boxIRI *url.URL) (crypto.PublicKey, httpsig.Algorithm, error) {
	if boxIRI != a.outboxURL {
		return nil, httpsig.RSA_SHA256, fmt.Errorf("unknown outbox url %s", boxIRI)
	} else if publicKeyId != a.actorURL.String() {
		return nil, httpsig.RSA_SHA256, fmt.Errorf("unknown public key id %q", publicKeyId)
	}
	return a.pubKey, httpsig.RSA_SHA256, nil
}

func (a *app) OnFollow(c context.Context, s *streams.Follow) pub.FollowResponse {
	return pub.AutomaticAccept
}

func (a *app) Unblocked(c context.Context, actorIRIs []*url.URL) error {
	return nil
}

func (a *app) FilterForwarding(c context.Context, activity vocab.ActivityType, iris []*url.URL) ([]*url.URL, error) {
	// Do NOT do this in real implementations. This turns the server into a
	// spambot. See the documentation in go-fed/activity/pub.
	return iris, nil
}

func (a *app) NewSigner() (httpsig.Signer, error) {
	s, _, err := httpsig.NewSigner([]httpsig.Algorithm{httpsig.RSA_SHA256}, nil, httpsig.Signature)
	return s, err
}

func (a *app) PrivateKey(boxIRI *url.URL) (privKey crypto.PrivateKey, pubKeyId string, err error) {
	return a.privKey, a.actorURL.String(), nil
}
