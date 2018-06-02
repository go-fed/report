package main

import (
	"crypto/tls"
	"flag"
	"github.com/go-fed/report"
	"net/http"
)

const (
	httpScheme  = "http"
	httpsScheme = "https"
)

var https *bool = flag.Bool("https", false, "enable serving via https")
var host *string = flag.String("host", "", "host domain of the server")
var newPath *string = flag.String("newPath", "/new", "path to newly created items")
var certFile *string = flag.String("cert", "", "tls cert file")
var keyFile *string = flag.String("key", "", "tls key file")

func main() {
	// Flags
	flag.Parse()
	scheme := httpScheme
	if *https {
		scheme = httpsScheme
	}

	// Server set up
	mux := http.NewServeMux()
	err := report.SetReportMux(mux, scheme, *host, *newPath)
	if err != nil {
		panic(err)
	}
	s := &http.Server{
		Addr:    ":" + scheme,
		Handler: mux,
	}

	// Run the server
	if *https {
		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
		s.TLSConfig = tlsConfig
		s.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
		if err := s.ListenAndServeTLS(*certFile, *keyFile); err != http.ErrServerClosed {
			panic(err)
		}
	} else {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}
}
