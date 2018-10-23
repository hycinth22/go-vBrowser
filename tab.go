package vBrowser

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

type Tab struct {
	Location     *url.URL
	PageResponse *http.Response
	// Document     *goquery.Document
	browser        *Browser
	httpClient     *http.Client
	ajaxHttpClient *http.Client
	BeforeAjax     []func(r *http.Request)
	AfterAjax      []func(r *http.Response)
}

// Load a new page in the tab (without referer)
func (t *Tab) Open(url string) (err error) {
	return t.load(url, false)
}

// jump to a page in the tab (with referer)
func (t *Tab) Jump(jumpUrl string) (err error) {
	return t.load(jumpUrl, true)
}

func (t *Tab) load(loadUrl string, withReferer bool) error {
	if t.PageResponse != nil {
		t.PageResponse.Body.Close()
	}
	location, err := url.Parse(loadUrl)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, loadUrl, nil)
	if err != nil {
		return errors.New("browser:" + err.Error())
	}

	if withReferer {
		// log.Println("set Referer", loadUrl)
		req.Header.Set("Referer", loadUrl)
	}
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return errors.New("browser:" + err.Error())
	}

	/* doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return errors.New("browser:" + err.Error())
	}*/

	t.Location = location
	t.PageResponse = resp
	// t.Document = doc
	return nil
}

// request a url (ajax)
func (t *Tab) Request(method, url string, body io.Reader) (page *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New("browser:" + err.Error())
	}
	return t.DoAjax(req)
}

// ajax
func (t *Tab) DoAjax(req *http.Request) (page *http.Response, err error) {
	if req.URL.Host != t.Location.Host || req.URL.Scheme != t.Location.Scheme {
		req.Header.Set("Origin", t.Location.Scheme+"://"+t.Location.Host)
	}
	req.Header.Set("Referer", t.Location.String())
	req.Header.Set("X-requested-with", "XMLHttpRequest")
	return t.ajaxHttpClient.Do(req)
}

func (t *Tab) GetPageCookie() []*http.Cookie {
	return t.browser.CookieJar.Cookies(t.Location)
}
