package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
	"bufio"
	"net"
	"runtime"
//	"bytes"
//	"io"
)

func listenConn() {						//Connection listener
	srv,err := net.Listen("tcp",":7777")
	defer srv.Close()
	fmt.Println("START")
	if err!=nil {
		fmt.Println("[CONN-LISTEN-FAIL] "+err.Error())
		os.Exit(1)	//Shmeh
	}
	for {
		conn,errr := srv.Accept()
		if errr!=nil {
			fmt.Println("[CONN-INIT-FAIL] "+errr.Error())
		}
		go connProc(conn)
		//fmt.Println("Started proc")
	}
}

func connProc(c net.Conn) {
	defer c.Close()
	url,err := bufio.NewReader(c).ReadString('\n')		//Get URL to process
	url = "https://en.wikipedia.org"+strings.Replace(url,"\n","",-1)
	if err!=nil {
		fmt.Println("[CONN-READ-FAIL] "+err.Error())
		return
	}
	crawl(url,c)
	//_,errr := fmt.Fprintln(c,links)	//Respond with URLs
	/*if errr!=nil {
		fmt.Println("[CONN-RESPOND-FAIL]["+url+"] "+errr.Error())
		return
	}*/
	_,errrr := fmt.Fprintln(c,"complete")
	if errrr!=nil {
		fmt.Println("[CONN-RESPOND-FAIL]["+url+"] "+errrr.Error())
		return
	}
	//fmt.Println("Completed proc")
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Extract all /wiki/ links from a given webpage
func crawl(url string,c net.Conn) {
//	client := &http.Client{}	//Create client
//	client.header.Set(
	resp, err := http.Get(url)
	//links := []string{}
//	fmt.Println(string(resp.Body[:]))
//	fmt.Println(err)

//	defer func() {
		// Notify that we're done after this function
		//out = strings.Join(links,"\n")
//	}()

	if err != nil {
		fmt.Println("[CRAWL-FAILURE] \"" + url + "\" "+err.Error())
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns
//	buf := new(bytes.Buffer)
//	buf.ReadFrom(b)
//	newStr := buf.String()

//	fmt.Printf(newStr)

//	fmt.Println(resp.Body.ReadAll())
	z := html.NewTokenizer(b)
	deduptable := make(map[string]bool)
	//wordtable := make(map[string]int)

	for {
		tt := z.Next()
//		fmt.Println(tt)
		switch tt {
		/*case html.TextToken:
			for _,w := range strings.Fields(string(z.Text()[:])) {
				wordtable[w]++
			}*/
		case html.ErrorToken:
			// End of the document, we're done
			return
//		case html.TextToken:
//			fmt.Println(string(z.Text()[:]))
		case html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// Make sure the url begines in /wiki/
			isWiki := strings.Index(url, "/wiki/") == 0
			if isWiki {
				if !deduptable[url] {	//Make sure we don't give credit twice
					fmt.Fprintln(c,url)
					deduptable[url]=true
				}
			}
		}
	}
}

/*func completionHandler(chUrls chan string) {
	// Channels
//	chUrls := make(chan string)

	for true {
		fmt.Println(<-chUrls)
		fmt.Println("complete")
	}
}*/

func main() {
/*	chUrls := make(chan string)
	go completionHandler(chUrls)		//Start completion handler
	bio := bufio.NewReader(os.Stdin)
	for true {
		url,_ := bio.ReadString('\n')
		url = strings.Replace(url,"\n","",-1)
		if url=="terminate" {
			os.Exit(0)
		}else{
			go crawl("https://en.wikipedia.org"+url,chUrls)
		}
	}*/
	runtime.GOMAXPROCS(runtime.NumCPU())
	listenConn()
}
