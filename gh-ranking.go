/*
gh-ranking is a simple tool to get the position 
of a programming language in popularity GitHub ranking.
*/
package main

import (
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go.net/html"
	"fmt"
	gnv "github.com/fern4lvarez/gonverter"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const base = "https://github.com/languages/"

// httpResponse type, enhanced http response
type httpResponse struct {
	url      string
	response *http.Response
	err      error
}

// timeoutDialer sets a connection with a given timeout
func timeoutDialer(secs int) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		c, err := net.Dial(netw, addr)
		if err != nil {
			return nil, err
		}
		c.SetDeadline(time.Now().Add(time.Duration(secs) * time.Second))
		return c, nil
	}
}

// getResp returns the httpResponse of the requested language
// using a a goroutine and a channel
func getResp(base string, lang string) *httpResponse {
	url := base + lang
	cl := http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialer(5),
		},
	}
	ch := make(chan *httpResponse)
	var response *httpResponse
	go func(url string) {
		resp, err := cl.Get(url)
		ch <- &httpResponse{url, resp, err}
	}(url)
	for {
		select {
		case response := <-ch:
			return response
		}
	}
	return response
}

// posRegexp extracts a decimal number from the given string
func posRegexp(s string) string {
	re := regexp.MustCompile(`\d+`)
	return re.FindString(s)
}

// encode a string to make it URL compatible
func encode(s string) string {
	return strings.Replace(s, " ", "%20", -1)
}

// Position returns the position on GitHub of the requested language
func Position(lang string) (int, error) {
	// Anon function that extracts the needed information from the tree node
	var p string
	f := func(n *html.Node) {
		if strings.Contains(n.Data, "the most") {
			p = "1"
			return
		} else if strings.Contains(n.Data, "most popular language on GitHub") {
			p = n.Data
			return
		}
	}

	// fetch the response
	r := getResp(base, lang)
	defer r.response.Body.Close()
	if r.err != nil || r.response.StatusCode != 200 {
		return 0, r.err
	}

	// Creates a tree from the response body
	tree, err := h5.New(r.response.Body)
	if err != nil {
		return 0, err
	}

	// Walk nodes of the tree until finds the needed data 
	// and converts it to an integer
	tree.Walk(f)
	position := gnv.StoI((posRegexp(p)))
	return position, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Hey, gimme just one language please!")
		os.Exit(1)
	}

	lang := os.Args[1]
	position, err := Position(lang)
	if position == 0 || err != nil {
		fmt.Printf("Something went wrong! Try again with some other real language, Mr. @HipsterHacker.\n")
		os.Exit(1)
	}

	fmt.Printf("%s is on position #%d on GitHub.\n", lang, position)
	fmt.Printf("For more information, visit %s%s\n", base, encode(lang))
}
