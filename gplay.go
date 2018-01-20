package main

import (
	"io"
	"fmt"
	"strings"
	"io/ioutil"
	"net/url"
	"net/http"
	"code.google.com/p/go.net/html"
)


type AppInfo struct {
	name string
	packageName string
}

// page: 1,2,...
func fetchPage(query string, page int, country string) (string, error) {
	start := (page - 1) * 24 // 0...
	num := 24 // google play's site doesn't accept >24
	escapedKeyword := url.QueryEscape(query)
	url := fmt.Sprintf("https://play.google.com/store/search?q=%s&c=apps&sort=%d&start=%d&num=%d&hl=%s",
		escapedKeyword, 
		1,  // 0:popularity, 1:relevence
		start,
		num,
		country)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func parseHtml(r io.Reader) (ret []AppInfo) {
	d := html.NewTokenizer(r)
	for { 
		tt := d.Next()
		if tt == html.ErrorToken {
			//fmt.Printf("err token\n")
			return
		}
		token := d.Token()
		switch tt {
			case html.StartTagToken:
				if (getAttr(token.Attr, "class") == "title" && 
					hasAttr(token.Attr, "title") && 
					hasAttr(token.Attr, "href")) {
					url, _ := url.Parse(getAttr(token.Attr, "href"))
					id := url.Query().Get("id")
					info := AppInfo{ name:getAttr(token.Attr, "title"), packageName:id }
					ret = append(ret, info)
				 }
			case html.TextToken:
			case html.EndTagToken:
			case html.SelfClosingTagToken:

		}
	}
	return ret
}

func hasAttr(attrs []html.Attribute, key string) bool {
	for _, attr := range(attrs) {
		if attr.Key == key {
			return true
		}
	}
	return false
}

func getAttr(attrs []html.Attribute, key string) string {
	for _, attr := range(attrs) {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getRanking(pkg, keyword, lang string) (int, error) {
	rank := 0
	for page:=1; page <= 20; page++ { // 20 seems to be the limit
		html, err := fetchPage(keyword, page, lang)
		if err != nil {
			return 0, err
		}
		ret := parseHtml(strings.NewReader(html))

		for _, info := range(ret) {
			rank++
			if pkg == info.packageName {
				return rank, nil
			}
		}
		if len(ret) < 24 {
			return 0, nil
		}
	}
	return 0, nil
}

