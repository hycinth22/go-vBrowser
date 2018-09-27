package vBrowser

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"

type Browser struct {
	UserAgent string
	CookieJar http.CookieJar
}

func NewBrowser() (browser *Browser) {
	var err error
	browser = new(Browser)
	browser.CookieJar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}
	browser.UserAgent = defaultUserAgent
	return
}

func (b *Browser) NewTab() (tab *Tab) {
	tab = new(Tab)
	tab.browser = b
	tab.httpClient = new(http.Client)
	tab.httpClient.Jar = b.CookieJar
	tab.httpClient.Transport = &transport{
		browser: b,
		tab:     tab,
	}
	tab.ajaxHttpClient = new(http.Client)
	tab.ajaxHttpClient.Jar = b.CookieJar
	tab.ajaxHttpClient.Transport = &ajaxTransport{
		transport{
			browser: b,
			tab:     tab,
		},
	}
	return
}
