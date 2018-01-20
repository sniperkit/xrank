package main

import (
	"fmt"
	"testing"
)

var msgFail = "%v function fails. Expects %v, returns %v"

func TestTimeoutDialer(t *testing.T) {
	secs := 20
	f1 := ExportTimeoutDialer(secs)
	c, err := f1("tcp", "google.com:80")
	if err != nil {
		t.Errorf(msgFail, "Timeout Dialer", "TCPConn", c)
	}
}

func TestGetResp(t *testing.T) {
	lang := "go"

	r := ExportGetResp(base, lang)
	if r.url != base+lang {
		t.Errorf(msgFail, "getResp: url", base+lang, r.url)
	} else if hr := fmt.Sprintf("%T", r.response); hr != "*http.Response" {
		t.Errorf(msgFail, "getResp: response", "*http.Response", hr)
	} else if r.err != nil {
		t.Errorf(msgFail, "getResp: err", "nil", r.err)
	}
}

func TestPosRegexp(t *testing.T) {
	s := "This language is the #23, not the best"
	if p := ExportPosRegexp(s); p != "23" {
		t.Errorf(msgFail, "posRegexp", "23", p)
	}

	s = "1"
	if p := ExportPosRegexp(s); p != "1" {
		t.Errorf(msgFail, "posRegexp", "1", p)
	}

	s = "No number here."
	if p := ExportPosRegexp(s); p != "" {
		t.Errorf(msgFail, "posRegexp", "", p)
	}
}

func TestEncode(t *testing.T) {
	s := "Ruby"
	if e := ExportEncode(s); e != "Ruby" {
		t.Errorf(msgFail, "encode", "Ruby", e)
	}

	s = "Visual Basic"
	if e := ExportEncode(s); e != "Visual%20Basic" {
		t.Errorf(msgFail, "encode", "Visual%20Basic", e)
	}
}

func TestPosition(t *testing.T) {
	// Ruby never will get down of the #3 position on GitHub lol
	lang := "Ruby"
	if p, err := Position(lang); p < 1 || p > 3 {
		t.Errorf(msgFail, "Position Ruby", "1, 2 or 3", p)
	} else if err != nil {
		t.Errorf(msgFail, "Position: error", "nil", err)
	}

	lang = "HipsterLanguage"
	if p, _ := Position(lang); p != 0 {
		t.Errorf(msgFail, "Position HipsterLanguage", "0", p)
	}
}

func ExamplePosition() {
	lang := "go"
	pos, err := Position(lang)
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println(pos)
	// Output: 23
}
