package apm

import (
	"net/http"

	"github.com/Toflex/directory_v2/pkg/apm/sentry"
	"github.com/go-resty/resty/v2"
)

func instrumentedRoundTripper() http.RoundTripper {
	return sentry.ClientOpt().HTTPTransport
}

func HTTPClientRequest() *resty.Request {
	return resty.New().SetTransport(instrumentedRoundTripper()).R()
}
