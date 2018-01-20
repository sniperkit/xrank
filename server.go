package main

import (
//	"os"
	"os/exec"
//	"io"
	"bufio"
	"fmt"
	"strings"
//	"time"
	"net"
//	"strconv"
)

type completion struct{
	links []string
	node string
}

func runJob(node string, url string, out chan *completion) {
	c,err := net.Dial("tcp",node+":7777")	//Connect
	o := []string{}
	defer c.Close()
	defer func() {
		//Create completion struct
		comp := new(completion)
		comp.links = o
		comp.node = node
		out <- comp
	}()
	if err!=nil {
		fmt.Println("[ERR][CONN-FAIL]["+node+"] "+err.Error())
		return
	}
	fmt.Fprintln(c,strings.Replace(url,"\n","",-1))			//Send url
	r := bufio.NewReader(c)
	kg := true
	for kg {
		s,_ := r.ReadString('\n')
		s = strings.Replace(s,"\n","",-1)
		//fmt.Println(s)
		if s=="complete" {
			fmt.Println("[Complete] "+s)
			kg = false
		}else{
			o = append(o,s)
		}
	}
}

func worker(ip string,start chan bool) {
	cf := exec.Command("ssh",ip,"go","run","worker.go")
	op,_ := cf.StdoutPipe()		//get output pipe
	reader := bufio.NewReader(op)
	//defer reader.Close()
//	defer cf.Close()
		fmt.Println("[START-WORKER]["+ip+"]")
		cf.Start()
		shmeh,_ := reader.ReadString('\n')
		//fmt.Println(err.Error())
		fmt.Println(shmeh)
		start <- true
		fmt.Println("[STARTED-WORKER]["+ip+"]")
	go func() {
		cf.Wait()
		fmt.Println("[TERMINATE]["+ip+"]")
	}()
	for {
		line,err := reader.ReadString('\n')		//Get line
		line = strings.Replace(line,"\n","",-1)
		fmt.Println("[STDOUT]["+ip+"] "+line)
		if err!=nil {
			fmt.Println("[ERR-READ]["+ip+"] "+err.Error());
			//os.Exit(1)
		}

	}
}
func main() {
	//Initalize everything
	linkcount := 0
	linkrank := make(map[string]int)
	completions := make(chan *completion,128)
	toindex := []string{"/wiki/Jay_Leno"}
	workerip := []string{"10.0.0.34","10.0.0.47","10.0.0.91","10.0.0.92","10.0.0.140","10.0.0.173"}
	sch := make(map[string]int)

	//Open xterm to display progress
	s := exec.Command("tee","stat")
	//spipe,_ := s.StdinPipe()
	s.Start()
	//stime := time.Now()
	fmt.Println(exec.Command("xterm","cat","stat").Start())	//Start an xterm to display output

	//Start Workers
	stchan := make(chan bool)
	for _,ip := range workerip {
		go worker(ip,stchan)
		fmt.Println(<-stchan)
		sch[ip]=0
	}
	fmt.Println("Started workers")

	//Dispatcher
	dispatch := func (url string) (bool) {
		//Loop through to find node that is not at full
		//work = false
		for nod,cnt := range sch {
			if cnt<65 {	//Found node not at full load
				go runJob(nod,url,completions)
				sch[nod]++
				return true
			}
		}
		return false
	}

	//Run indexing loop
	goal := 65536
	for (linkcount<goal) {
		/*if (linkcount%4)==0 {			//Print out progress
			//fmt.Println("Prog")
			speed := float64(linkcount)/time.Since(stime).Seconds()
			timeRemaining := float64(goal)/speed
			fmt.Fprintln(spipe,string(linkcount)+"l "+strconv.FormatFloat(speed,'f',2,64)+"l/s "+strconv.FormatFloat(timeRemaining,'f',2,64)+"s r")
			//fmt.Println("prog")
		}*/
		if (len(toindex)>0) && dispatch(toindex[0]) {	//Try to dispatch if there is smth to dispatch
			toindex = append(toindex[:0], toindex[1:]...)
		} else {
			comp := <-completions
			dwo := true
			for _,l := range comp.links {
				if linkrank[l]==0 {
					if dwo {	//Try to dispatch
						dwo = dispatch(l)
					}
					if !dwo {	//Didn't work or dispatch failed previously
						toindex = append(toindex,l)
					}
					linkrank[l]=1
				} else {
					linkrank[l]++
				}
			}
			sch[comp.node]--
			fmt.Println(sch[comp.node])
			linkcount++
		}
	}
}
