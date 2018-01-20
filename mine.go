package main

import (
	"fmt"
	"github.com/steaz/reddilytics/mining"
	"github.com/steaz/reddilytics/model"
	"log"
	"time"
)

func main() {
	err := model.SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	mining.UpdateAll()
	fmt.Println("Done: ", time.Now())
	fmt.Println("\n--------------------------------------\n")
}
