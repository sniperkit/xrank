package smartb

import (
	"fmt"
)

type holder interface {
	ToString() string
}

type Balancer interface {
	// getting elements
	Best() *element
	Worst() *element

	// "feed" or "punish" elements
	Good(*element)
	Bad(*element)

	set(*element)
}

type Container struct {
	balancer Balancer
	current  *element
	all      map[string]*element
}
type element struct {
	next, prev *element
	item       holder
}

func NewContainer(optimizer int, warming int) (*Container, error) {
	cnt := new(Container)

	bl, err := newBalancer(optimizer)
	if err != nil {
		return nil, err
	}
	cnt.balancer = bl
	cnt.all = make(map[string]*element)
	return cnt, nil
}

func (c *Container) Add(elem holder) {
	var key = elem.ToString()
	// TODO: mutex 2
	if obj, ok := c.all[key]; ok {
		fmt.Printf("Already in container: %p", obj.item)
	} else {
		el := new(element)
		el.item = elem
		c.all[key] = el
		c.balancer.set(el)
	}
}

func (c *Container) GetElem() *element {

	if len(c.all) == 0 {
		return nil
	}

	return c.balancer.Best()
}

func (e *element) Item() interface{} {
	return e.item
}

func (c *Container) Len() int {
	return len(c.all)
}
