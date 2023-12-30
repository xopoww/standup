package httputil

import (
	"context"
	"net/http"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func CancellableRoundTripper(stop <-chan struct{}, inner http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(r *http.Request) (rsp *http.Response, err error) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		r = r.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			rsp, err = inner.RoundTrip(r) //nolint:bodyclose // we actually return rsp here
			close(done)
		}()

		select {
		case <-stop:
			cancel() // cancel the request
			<-done   // wait when RoundTrip fails with "context cancelled"
		case <-done:
		}
		return rsp, err
	})
}
