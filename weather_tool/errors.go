package weathertool

import "errors"

// ErrUnknownLocation is returned when a province or county name does not
// match any entry in the allowlist (Provinces / Counties).
var ErrUnknownLocation = errors.New("weathertool: unknown location")

// UpstreamError wraps a failure talking to the met.ie XML feed: a transport
// error, a non-200 response, or a body that failed to unmarshal.
// StatusCode is 0 for transport/parse failures (no HTTP response received).
type UpstreamError struct {
	URL        string
	StatusCode int
	Err        error
}

func (e *UpstreamError) Error() string {
	if e.StatusCode != 0 {
		return "weathertool: upstream " + e.URL + ": " + e.Err.Error()
	}
	return "weathertool: upstream " + e.URL + ": " + e.Err.Error()
}

func (e *UpstreamError) Unwrap() error {
	return e.Err
}
