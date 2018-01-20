package smartb

import (
	"errors"
	"sync"
)

const (
	UBalancer = 1 + iota
)

type (
	uniformNode struct {
		sync.RWMutex
		counter int
	}
	uniformBalancer struct {
		sync.RWMutex
		ranks map[*element]*uniformNode
	}
)

func (un *uniformNode) get() int {
	un.RLock()
	defer un.RUnlock()
	return un.counter
}

func (un *uniformNode) inc() {
	un.Lock()
	un.counter++
	un.Unlock()
}

func (ub *uniformBalancer) set(el *element) {
	ub.Lock()
	if _, ok := ub.ranks[el]; !ok {
		ub.ranks[el] = new(uniformNode)
	}
	ub.Unlock()

}

func (ub *uniformBalancer) Best() *element {

	var max_el *element
	var max_node *uniformNode
	var max_val int

	ub.RLock()
	for k, val := range ub.ranks {
		if max_el == nil {
			max_el = k
			max_node = ub.ranks[k]
			max_val = val.get()
		} else {
			tmp := val.get()
			if tmp < max_val {
				max_el = k
				max_node = ub.ranks[k]
				max_val = tmp
			}
		}
	}
	ub.RUnlock()

	if max_node != nil {
		max_node.inc()
	}

	return max_el
}
func (ub *uniformBalancer) Worst() *element {
	return nil
}

func (ub *uniformBalancer) Good(el *element) {
}
func (ub *uniformBalancer) Bad(el *element) {
}

func newBalancer(btype int) (Balancer, error) {
	switch btype {
	case UBalancer:
		ub := new(uniformBalancer)
		ub.ranks = make(map[*element]*uniformNode)
		return ub, nil
	}
	return nil, errors.New("Invalid balancer type")
}
