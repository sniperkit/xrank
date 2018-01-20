package lib

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	//"net/http"
	"net/url"
)

type Request struct {
}

func (this *Request) Get(api string, params map[string]string, result interface{}) {
	agent := new()
	u, _ := url.Parse(api)
	q := u.Query()

	if params != nil {
		for key, value := range params {
			q.Set(key, value)
		}
	}
	u.RawQuery = q.Encode()

	fmt.Printf("[REQUEST]: %s \n", u.String())

	resp, body, errs := agent.Get(u.String()).End()
	if len(errs) > 0 {
		fmt.Println(errs)
		panic(errs[0])
	}
	err := json.Unmarshal([]byte(body), result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status code:", resp.StatusCode)
}

func new() *gorequest.SuperAgent {
	agent := gorequest.New()
	//agent.SetDebug(true)
	agent.Type("json")
	return agent
}
