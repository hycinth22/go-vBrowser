package vBrowser

import (
	"net/http"
)

type transport struct {
	browser *Browser
	tab     *Tab
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("User-Agent", t.browser.UserAgent)
	resp, err = http.DefaultTransport.RoundTrip(req)
	return
}

type ajaxTransport struct {
	transport
}

func (t *ajaxTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	for _, f:= range t.tab.BeforeAjax{
		f(req)
	}
	resp, err = t.transport.RoundTrip(req)
	for _, f:= range t.tab.AfterAjax{
		f(resp)
	}
	return resp, err
}
