// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

package acme

import (
	"net/http"
	"strconv"
	"time"

	"tailscale.com/tsweb"
)

// CertRateLimitedError is returned by [LocalBackend.GetCertPEMWithValidity]
// when the upstream ACME CA rate-limited the issuance. RetryAfter is the
// CA's suggested wait, or zero if none was provided.
type CertRateLimitedError struct {
	RetryAfter time.Duration
	Underlying error
}

func (e *CertRateLimitedError) Error() string { return e.Underlying.Error() }
func (e *CertRateLimitedError) Unwrap() error { return e.Underlying }

// HTTPStatus implements [tsweb.HTTPStatuser].
func (e *CertRateLimitedError) HTTPStatus() tsweb.HTTPError {
	h := http.Header{}
	if e.RetryAfter > 0 {
		h.Set("Retry-After", strconv.Itoa(int(e.RetryAfter.Seconds())))
	}
	return tsweb.HTTPError{
		Code:   http.StatusTooManyRequests,
		Msg:    e.Error(),
		Header: h,
	}
}
