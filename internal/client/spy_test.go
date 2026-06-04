package client_test

import (
	"io"
	"net/http"
)

type spyReadCloser struct {
	closes int
	reads  int
	real   io.ReadCloser
}

func (s *spyReadCloser) Read(p []byte) (n int, err error) {
	s.reads++
	return s.real.Read(p)
}

func (s *spyReadCloser) Close() error {
	s.closes++
	return s.real.Close()
}

type spyRoundTripper struct {
	rc   spyReadCloser
	real http.RoundTripper
}

func (rt *spyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := rt.real.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	rt.rc.real = resp.Body
	resp.Body = &rt.rc
	return resp, err
}
