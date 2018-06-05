package report

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/vocab"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const (
	actorPath     = "/actor"
	inboxPath     = "/actor/inbox"
	outboxPath    = "/actor/outbox"
	followingPath = "/actor/following"
	followersPath = "/actor/followers"
	likedPath     = "/actor/liked"
	authPath      = "/auth"
	tokenPath     = "/token"
)

// SetReportMux builds a basic Social API and Federate API server using the
// bare-bones go-fed/activity library. Due to the test suite, a skeleton
// SocialAPIVerifier is used to stub out the OAuth 2 calls but this server does
// NOT do any actual authentication. It authorizes everyone. Do not use this for
// any other purpose but generating an implementation report. It should NOT be
// used as a reference for building an actual implementation.
//
// You have been thoroughly warned.
func SetReportMux(m *http.ServeMux, scheme, host, newPath string) error {
	// Implementation specific data
	actorURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, actorPath))
	if err != nil {
		return err
	}
	inboxURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, inboxPath))
	if err != nil {
		return err
	}
	outboxURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, outboxPath))
	if err != nil {
		return err
	}
	followingURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, followingPath))
	if err != nil {
		return err
	}
	followersURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, followersPath))
	if err != nil {
		return err
	}
	likedURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, likedPath))
	if err != nil {
		return err
	}
	authURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, authPath))
	if err != nil {
		return err
	}
	tokenURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, tokenPath))
	if err != nil {
		return err
	}
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	pubKey := privKey.Public()
	endpoints := &vocab.Object{}
	endpoints.SetOauthAuthorizationEndpoint(authURL)
	endpoints.SetOauthTokenEndpoint(tokenURL)
	actor := &vocab.Person{}
	actor.SetEndpoints(endpoints)
	actor.SetId(actorURL)
	actor.AppendNameString("Implementation Report Account")
	actor.SetOutboxAnyURI(outboxURL)
	actor.SetInboxAnyURI(inboxURL)
	actor.SetFollowingAnyURI(followingURL)
	actor.SetFollowersAnyURI(followersURL)
	actor.SetLikedAnyURI(likedURL)
	actor.SetPreferredUsername("Implementation Report Account")

	// Prepare basic implementation
	verifier := &doNotUseThisItIsNotOAuth{
		ActorURL:  actorURL,
		OutboxURL: outboxURL,
	}
	app := newApp(scheme, host, newPath, actorURL, inboxURL, outboxURL, followingURL, followersURL, likedURL, pubKey, privKey, actor, verifier)
	fedCb := &nothingCallbacker{}
	socialCb := &nothingCallbacker{}
	clock := &localClock{}
	deliverer := &syncDeliverer{}
	pubber := pub.NewPubber(clock, app, socialCb, fedCb, deliverer, &http.Client{}, "go-fed-report", 5, 5)
	serveFn := pub.ServeActivityPubObject(app, clock)
	addMissingFn := func(r *http.Request) {
		r.URL.Host = host
		r.URL.Scheme = scheme
	}

	// Set up sync primitives
	lockKey := 1
	lockKeyMu := &sync.Mutex{}
	getLockKeySafely := func() (context.Context, context.CancelFunc) {
		c := context.Background()
		k := lockKeyType("lockKey")
		lockKeyMu.Lock()
		defer lockKeyMu.Unlock()
		v := lockKey
		lockKey++
		return context.WithCancel(context.WithValue(c, k, v))
	}
	// Set up handlers
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		addMissingFn(r)
		log.Printf("received request to %q", r.URL)
		c, cfn := getLockKeySafely()
		defer cfn()
		if handled, err := serveFn(c, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		} else if handled {
			return
		}
		log.Printf("request to %q not an activitypub request", r.URL)
		log.Print(r)
		w.WriteHeader(http.StatusNotFound)
	})
	m.HandleFunc(actorPath, func(w http.ResponseWriter, r *http.Request) {
		addMissingFn(r)
		log.Printf("received request to %q", r.URL)
		c, cfn := getLockKeySafely()
		defer cfn()
		if handled, err := serveFn(c, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		} else if handled {
			return
		}
		log.Printf("request to %q not an activitypub request", r.URL)
		log.Print(r)
		w.WriteHeader(http.StatusNotFound)
	})
	m.HandleFunc(inboxPath, func(w http.ResponseWriter, r *http.Request) {
		addMissingFn(r)
		log.Printf("received request to %q", r.URL)
		c, cfn := getLockKeySafely()
		defer cfn()
		if handled, err := pubber.GetInbox(c, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		} else if handled {
			return
		}
		if handled, err := pubber.PostInbox(c, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		} else if handled {
			return
		}
		log.Printf("request to %q not an activitypub request", r.URL)
		log.Print(r)
		w.WriteHeader(http.StatusNotFound)
	})
	m.HandleFunc(outboxPath, func(w http.ResponseWriter, r *http.Request) {
		addMissingFn(r)
		log.Printf("received request to %q", r.URL)
		c, cfn := getLockKeySafely()
		defer cfn()
		if handled, err := pubber.GetOutbox(c, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		} else if handled {
			return
		}
		if handled, err := pubber.PostOutbox(c, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		} else if handled {
			return
		}
		log.Printf("request to %q not an activitypub request", r.URL)
		log.Print(r)
		w.WriteHeader(http.StatusNotFound)
	})
	m.HandleFunc(authPath, func(w http.ResponseWriter, r *http.Request) {
		addMissingFn(r)
		log.Printf("received request to %q", r.URL)
		verifier.AuthorizeRequestWithoutActuallyDoingAnything(w, r)
	})
	m.HandleFunc(tokenPath, func(w http.ResponseWriter, r *http.Request) {
		addMissingFn(r)
		log.Printf("received request to %q", r.URL)
		verifier.GrantBearerTokenWithoutActuallyDoingAnything(w, r)
	})
	return nil
}
