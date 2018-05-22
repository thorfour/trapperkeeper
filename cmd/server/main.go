package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/acme/autocert"

	"github.com/thorfour/trapperkeeper/pkg/keeper"
)

const (
	ephemeral = "ephemeral"
	inchannel = "in_channel"
	dataDir   = "."
)

var (
	port  = flag.Int("p", 443, "port to serve on")
	debug = flag.Bool("d", false, "turn TLS off")
	// AllowedHost is the ACME allowed host
	AllowedHost string
	// SupportEmail email for ACME provider to contact for TLS problems
	SupportEmail string
	// RedisAddr is the redis enpoint to be used
	RedisAddr string
	// RedisPw is if the endpoing has a password
	RedisPw string
)

// response is the json struct for a slack response
type response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func init() {
	flag.StringVar(&AllowedHost, "host", "api.stocktopus.io", "ACME allowed FQDN")
	flag.StringVar(&SupportEmail, "email", "support@stocktopus.io", "ACME support email")
	RedisAddr = os.Getenv("REDISADDR")
	RedisPw = os.Getenv("REDISPW")
}

func main() {
	flag.Parse()
	keeper.RedisAddr = RedisAddr
	keeper.RedisPw = RedisPw
	log.Printf("%s", AllowedHost)
	log.Printf("Starting server on port %v", *port)
	run(*port, *debug)
}

func run(p int, d bool) {
	if d {
		http.HandleFunc("/v1", handler)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", p), nil))
	} else {
		mux := &http.ServeMux{}
		mux.HandleFunc("/v1", handler)
		hostPolicy := func(ctx context.Context, host string) error {
			if host == AllowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s allowed", AllowedHost)
		}
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
			Email:      SupportEmail,
		}
		srv := &http.Server{
			Handler: mux,
			Addr:    fmt.Sprintf(":%v", p),
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
			},
		}
		go http.ListenAndServe(":80", m.HTTPHandler(nil))
		log.Fatal(srv.ListenAndServeTLS("", ""))
	}
}

func handler(resp http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	text := strings.Split(req.Form["text"][0], " ")
	msg, err := keeper.Handle(text[0], req.Form["user_id"][0], req.Form["user_name"][0], text[1:])
	newReponse(resp, msg, err)
}

func newReponse(resp http.ResponseWriter, message string, err error) {
	r := &response{
		ResponseType: inchannel,
		Text:         message,
	}

	// Swithc to an ephemeral message
	if err != nil {
		r.ResponseType = ephemeral
		r.Text = err.Error()
	}

	b, err := json.Marshal(r)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(b)
	return
}

func hostPolicy(ctx context.Context, host string) error {
	if host == AllowedHost {
		return nil
	}

	return fmt.Errorf("acme/autocert: only %s hist is allowed", AllowedHost)
}
