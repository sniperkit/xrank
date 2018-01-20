package main

import (
	"fmt"
	"github.com/steaz/reddilytics/analytics"
	"github.com/steaz/reddilytics/model"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	err := model.SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	preFile, err := os.Open("pretext.md")
	if err != nil {
		log.Fatal(err)
	}
	pretext, err := ioutil.ReadAll(preFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(pretext))
	analytics.Analyze()
}
