package vBrowser

import (
	"net/http"
	"testing"
)

const testURL = "http://www.msftconnecttest.com/redirect"

var tab = b.NewTab()

type testTransport struct {
	actual   http.RoundTripper
	testFunc func(req *http.Request, resp *http.Response)
}

func (t *testTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.actual.RoundTrip(req)
	if t.testFunc != nil {
		t.testFunc(req, resp)
	}
	return resp, err
}

type testFunc func(req *http.Request, resp *http.Response)

func hookTransport(client *http.Client, f testFunc) {
	client.Transport = &testTransport{
		actual:   tab.httpClient.Transport,
		testFunc: f,
	}
}

func unhookTransport(client *http.Client) {
	if trans, ok := client.Transport.(*testTransport); ok {
		tab.httpClient.Transport = trans.actual
	}
}

func TestTab_Open(t *testing.T) {
	// hook the transport
	var ok bool
	c := tab.httpClient.CheckRedirect
	tab.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	} // stop auto redirect(because it cause referer change)
	testFunc := func(req *http.Request, resp *http.Response) {
		referer := req.Header.Get("Referer")
		t.Log(referer)
		ok = referer == ""
	}
	hookTransport(tab.httpClient, testFunc)
	hookTransport(tab.ajaxHttpClient, testFunc)

	// test
	err := tab.Open(testURL)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !ok {
		t.Error("Test Referer Failed")
		t.Fail()
	}

	// resume the transport
	tab.httpClient.CheckRedirect = c
	unhookTransport(tab.httpClient)
	unhookTransport(tab.ajaxHttpClient)
}

func TestTab_Jump(t *testing.T) {
	// hook the transport
	var ok bool
	lastURL := tab.Location
	t.Log("lastURL:", lastURL)
	c := tab.httpClient.CheckRedirect
	tab.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	} // stop auto redirect(because it cause referer change)
	testFunc := func(req *http.Request, resp *http.Response) {
		referer := req.Header.Get("Referer")
		t.Log(referer)
		ok = referer == lastURL.String()
	}
	hookTransport(tab.httpClient, testFunc)
	hookTransport(tab.ajaxHttpClient, testFunc)

	// test
	err := tab.Jump(testURL)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !ok {
		t.Error("Test Referer Failed")
		t.Fail()
	}

	// resume the transport
	tab.httpClient.CheckRedirect = c
	if trans, ok := tab.httpClient.Transport.(*testTransport); ok {
		tab.httpClient.Transport = trans.actual
	}
	unhookTransport(tab.httpClient)
	unhookTransport(tab.ajaxHttpClient)
}

func TestTab_Request(t *testing.T) {
	_, err := tab.Request("1", "1", nil)
	if err == nil || err.Error() != `1 1: unsupported protocol scheme ""` {
		t.Log(err)
		t.Fail()
	}
	_, err = tab.Request("1", "http://1.1", nil)
	if err == nil || err.Error() != `1 http://1.1: dial tcp: lookup 1.1: no such host` {
		t.Log(err)
		t.Fail()
	}

	// hook the transport
	var ok bool
	c := tab.httpClient.CheckRedirect
	tab.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	} // stop auto redirect(because it cause referer change)

	testFunc := func(req *http.Request, resp *http.Response) {
		referer := req.Header.Get("Referer")
		method := req.Method
		url := req.URL.String()
		t.Log(method, url, referer)
		t.Log("referer", referer)
		t.Log("method", method)
		t.Log("url", url)
		ok = referer == tab.Location.String() && method == http.MethodPost && url == testURL
	}
	hookTransport(tab.httpClient, testFunc)
	hookTransport(tab.ajaxHttpClient, testFunc)

	// test
	resp, err := tab.Request(http.MethodPost, testURL, nil)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !ok {
		t.Error("not ok")
		t.Fail()
	}

	// resume the transport
	tab.httpClient.CheckRedirect = c
	unhookTransport(tab.httpClient)
	unhookTransport(tab.ajaxHttpClient)
}
