package vBrowser

import "testing"


var b = NewBrowser()

func TestDefaultUserAgent(t *testing.T) {
	if b.UserAgent != defaultUserAgent {
		t.FailNow()
	}
}
