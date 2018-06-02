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
)

const (
	actorPath  = "/actor"
	inboxPath  = "/actor/inbox"
	outboxPath = "/actor/outbox"
	authPath   = "/auth"
	tokenPath  = "/token"
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
	actor := &vocab.Object{}
	actor.SetEndpoints(endpoints)
	actor.SetId(actorURL)
	actor.AppendNameString("Implementation Report Account")
	actor.SetOutboxOrderedCollection(&vocab.OrderedCollection{})
	actor.SetInboxOrderedCollection(&vocab.OrderedCollection{})
	actor.SetFollowingCollection(&vocab.Collection{})
	actor.SetFollowersCollection(&vocab.Collection{})
	actor.SetLikedCollection(&vocab.Collection{})
	actor.SetPreferredUsername("Implementation Report Account")

	// Prepare basic implementation
	verifier := &doNotUseThisItIsNotOAuth{
		ActorURL:  actorURL,
		OutboxURL: outboxURL,
	}
	app := newApp(scheme, host, newPath, actorURL, inboxURL, outboxURL, pubKey, privKey, actor, verifier)
	fedCb := &nothingCallbacker{}
	socialCb := &nothingCallbacker{}
	clock := &localClock{}
	deliverer := &syncDeliverer{}
	pubber := pub.NewPubber(clock, app, socialCb, fedCb, deliverer, &http.Client{}, "go-fed-report", 5, 5)
	serveFn := pub.ServeActivityPubObject(app, clock)

	// Set up handlers
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request to %q", r.URL)
		c := context.Background()
		if handled, err := serveFn(c, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		} else if handled {
			return
		}
		log.Printf("request to %q not an activitypub request", r.URL)
		w.WriteHeader(http.StatusNotFound)
	})
	m.HandleFunc(actorPath, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request to %q", r.URL)
		c := context.Background()
		if handled, err := serveFn(c, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		} else if handled {
			return
		}
		log.Printf("request to %q not an activitypub request", r.URL)
		w.WriteHeader(http.StatusNotFound)
	})
	m.HandleFunc(inboxPath, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request to %q", r.URL)
		c := context.Background()
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
		w.WriteHeader(http.StatusNotFound)
	})
	m.HandleFunc(outboxPath, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request to %q", r.URL)
		c := context.Background()
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
		w.WriteHeader(http.StatusNotFound)
	})
	m.HandleFunc(authPath, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request to %q", r.URL)
		verifier.AuthorizeRequestWithoutActuallyDoingAnything(w, r)
	})
	m.HandleFunc(tokenPath, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request to %q", r.URL)
		verifier.GrantBearerTokenWithoutActuallyDoingAnything(w, r)
	})
	return nil
}
