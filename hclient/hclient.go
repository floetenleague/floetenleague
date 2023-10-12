package hclient

import "net/http"

var Client = &http.Client{Transport: &addUserAgent{}}

type addUserAgent struct{}

func (t *addUserAgent) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "OAuth floetenleague/0.0.0 (contact: EMAIL)")
	return http.DefaultTransport.RoundTrip(req)
}
