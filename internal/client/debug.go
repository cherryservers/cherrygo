package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
)

type debugger struct {
	wrapped requestDoer
	dst     io.Writer
	l       *log.Logger

	once sync.Once
}

func (d *debugger) Do(req *http.Request) (*http.Response, error) {
	d.once.Do(func() {
		d.l = log.New(d.dst, "cherrygo ", log.Default().Flags())
	})

	// Mask token.
	auth := req.Header.Get("Authorization")
	if auth != "" {
		req.Header.Set("Authorization", "***")
	}

	d.l.Printf("\nAPI Endpoint: %v\n", req.URL)

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	d.l.Printf("\n+++++++++++++REQUEST+++++++++++++\n%s\n+++++++++++++++++++++++++++++++++", string(dump))

	if auth != "" {
		req.Header.Set("Authorization", auth)
	}

	resp, err := d.wrapped.Do(req)
	if err != nil {
		return nil, err
	}

	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		// RoundTrip must always return err=nil if it obtained a response.
		dump = fmt.Appendf(nil, "failed to dump response: %v", err)
	}

	d.l.Printf("\n+++++++++++++RESPONSE+++++++++++++\n%s\n+++++++++++++++++++++++++++++++++", string(dump))
	return resp, nil
}
