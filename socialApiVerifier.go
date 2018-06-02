package report

import (
	"encoding/json"
	"fmt"
	"github.com/go-fed/activity/pub"
	"log"
	"net/http"
	"net/url"
)

var _ pub.SocialAPIVerifier = &doNotUseThisItIsNotOAuth{}

// Do not use this. Do not look at this. In fact, delete this. Provides zero
// security functionality. Mocks out the OAuth process. Does not actually do any
// authenticating. Authorizes everybody. Do not use this. Do not use this as a
// reference for actually implementing OAuth. I feel terrible for even writing
// this code.
type doNotUseThisItIsNotOAuth struct {
	ActorURL  *url.URL
	OutboxURL *url.URL
}

// Do not do this in real implementations. This does no actual authentication.
// Do not do this in real implementation.
func (o *doNotUseThisItIsNotOAuth) AuthorizeRequestWithoutActuallyDoingAnything(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	var redir *url.URL
	if rv, ok := v["redirect_uri"]; ok && len(rv) > 0 {
		rs, err := url.QueryUnescape(rv[0])
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		redir, err = url.Parse(rs)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if state, ok := v["state"]; ok && len(state) > 0 {
		for _, s := range state {
			redir.Query().Add("state", s)
		}
	}
	redir.Query().Add("code", "doNotDoThisInRealImplementations")
	w.Header().Set("Location", redir.String())
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusFound)
}

func (o *doNotUseThisItIsNotOAuth) GrantBearerTokenWithoutActuallyDoingAnything(w http.ResponseWriter, r *http.Request) {
	token := struct {
		A string `json:"access_token"`
		T string `json:"token_type"`
	}{
		A: "doNotDoThisInRealImplementations",
		T: "Bearer",
	}
	b, err := json.Marshal(token)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Write(b)
}

// Verify implements SocialAPIVerifier. Since this is a mock implementation,
// any and every one will be verified. Don't use this in real implementations.
func (o *doNotUseThisItIsNotOAuth) Verify(r *http.Request) (authenticatedUser *url.URL, authn, authz bool, err error) {
	bearer := r.Header.Get("Authorization")
	if bearer != "Bearer doNotDoThisInRealImplementations" {
		return nil, false, false, fmt.Errorf("bad bearer %q", bearer)
	}
	return o.ActorURL, true, true, nil
}

// VerifyForOutbox implements SocialAPIVerifier. Since this is a mock
// implementation, any and every one will be verified. Don't use this in real
// implementations.
func (o *doNotUseThisItIsNotOAuth) VerifyForOutbox(r *http.Request, outbox *url.URL) (authn, authz bool, err error) {
	if *outbox != *o.OutboxURL {
		err = fmt.Errorf("bad outbox url %q", outbox)
		return
	}
	_, authn, authz, err = o.Verify(r)
	return
}
