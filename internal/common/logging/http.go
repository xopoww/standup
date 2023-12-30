package logging

import (
	"net/http"
	"time"

	"github.com/xopoww/standup/internal/common/httputil"
	"go.uber.org/zap"
)

func RoundTripper(l *zap.Logger, inner http.RoundTripper) http.RoundTripper {
	s := l.Sugar()
	return httputil.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		s.Debugf("HTTP Req: %s %s.", r.Method, r.URL)

		start := time.Now()
		rsp, err := inner.RoundTrip(r)
		delta := time.Since(start)

		switch {
		case err != nil:
			s.Errorf("HTTP Err: %s %s (%s): %s.", r.Method, r.URL, delta, err)
		case rsp.StatusCode != http.StatusOK:
			s.Warnf("HTTP Rsp: %s %s (%s): status code %d.", r.Method, r.URL, delta, rsp.StatusCode)
		default:
			s.Debugf("HTTP Rsp: %s %s (%s): OK.", r.Method, r.URL, delta)
		}

		return rsp, err
	})
}
