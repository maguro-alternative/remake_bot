package mock

import (
	"net/http"
)

type roundTripFn func(req *http.Request) *http.Response

func (f roundTripFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewStubHttpClient(fn roundTripFn) *http.Client {
	return &http.Client{
		Transport: roundTripFn(fn),
	}
}
