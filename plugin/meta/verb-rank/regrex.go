package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"regexp"
	"time"
)

func main() {
	re := regexp.MustCompile("a(x*)b")
	fmt.Println(re.ReplaceAllString("-ab-axxb-", "T"))
	fmt.Println(re.ReplaceAllString("-ab-axxb-", "$1"))
	fmt.Println(re.ReplaceAllString("-ab-axxb-", "$1W"))
	fmt.Println(re.ReplaceAllString("-ab-axxb-", "${1}W"))
	length := len("a")
	fmt.Println(length)

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
	}
	defer c.Close()
	reply, err := c.Do("get", "mike")
	fmt.Println(reply, err)
	//channel()
}

func channel() {
	done := make(chan bool)
	go parse(done)
	product(done)
}

func product(done chan<- bool) {
	fmt.Println("start sleep 2 sec...")
	time.Sleep(2 * time.Second)
	fmt.Println("wake up...")
	done <- true
	time.Sleep(1 * time.Second)
}

func parse(done <-chan bool) {
	var finished bool
	for {
		if finished {
			break
		}
		select {
		case <-done:
			fmt.Println("get from chan done..")
			finished = true
		}
	}
	fmt.Println("finished...")
}
