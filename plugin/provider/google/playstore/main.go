package main

import (
	"fmt"
	"flag"
	"os"
)

var pkg string
var lang string

func init() {
	flag.StringVar(&pkg, "pkg", "", "package name")
	flag.StringVar(&lang, "lang", "en", "language code such as 'en'")
}

func main() {

	flag.Parse()

	if pkg == "" {
		fmt.Println("please specify -pkg <name>")
		os.Exit(1)
	}

	for _, keyword := range flag.Args() {
		rank, _ := getRanking(pkg, keyword, lang)
		fmt.Printf("#%d for %s\n", rank, keyword)
	}


}

