package lib

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"
)

const (
	CLONE_PATH = "repos"
)

type Git struct {
	Id          int64  `json:"id"`
	Owner       Owner  `json:"owner"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Private     bool   `json:"private"`
	HtmlUrl     string `json:"html_url"`
	Description string `json:"description"`
	GitUrl      string `json:"git_url"`
	CloneUrl    string `json:"clone_url"`
}

type Owner struct {
	Login string
	id    int64
}

func (this *Git) Clone() error {
	path := CLONE_PATH + "/" + this.FullName
	cmd := exec.Command("git", "clone", this.CloneUrl, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	log.Printf("git clone %s %s \n", this.CloneUrl, path)
	err := cmd.Run()
	if err != nil {
		log.Println("Error:", err)
		log.Printf("rm -rf %s \n", CLONE_PATH+"/"+this.Owner.Login)
		Execute("rm", "-rf", path)
		return err
	}
	fmt.Printf("Finished git clone: %s, output:%q\n", this.FullName, out.String())
	return nil
}

func Execute(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func Start(size int, done chan<- bool) {
	capaticy := size
	var lock sync.Mutex
	for {
		if capaticy <= 0 {
			break
		}
		thread := ThreadPool.Poll()
		if thread != nil {
			go func(capaticy *int, locker sync.Mutex) {
				repo := Queue.Poll()
				if repo != nil {
					git, ok := repo.(Git)
					if !ok {
						log.Fatalln("Get git from queue failed....")
						return
					}
					err := git.Clone()
					//After finished, re-push to ThreadPool
					ThreadPool.Push(1)
					if err == nil {
						log.Printf("Add path to repo queue: %s \n", git.FullName)
						RepoQueue.Push(git.FullName)
					}

					locker.Lock()
					(*capaticy)--
					locker.Unlock()
				}
			}(&capaticy, lock)
		}
		//sleep 1 seconds
		time.Sleep(time.Second * 1)
	}
	done <- true
}
