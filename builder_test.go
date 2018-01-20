package smartb

import "testing"
import "fmt"

// func TestSomething(t *testing.T) {
// 	r := new(Ring)
// 	t.Logf("Item: %v", r.Get())
// }

// func TestOrder(t *testing.T) {
// 	r := new(Ring)

// 	for i := 0; i < 10; i++ {

// 		testing_str := fmt.Sprintf("hi test me %d", i)
// 		new_elem := element{nil, nil, i, &testing_str}
// 		r.NewElem(&new_elem)
// 	}

// 	head := r.Head()
// 	for head != nil {
// 		itm := head.Item()
// 		if itm != nil {
// 			t.Logf("Item: %p <- (%p %s) -> %p", head.Prev(), head, *itm.(*string), head.Next())
// 		} else {
// 			t.Logf("Item: %+v", head.Item())
// 		}
// 		head = head.Next()
// 	}
// }

type Link string

func (l Link) ToString() string {
	return string(l)
}

func TestContainer(t *testing.T) {

	const elems_len = 10

	cont, err := NewContainer( /*optimize strategy*/ 1 /*warmup cache*/, 100)

	if err != nil {
		t.Fatalf("Error creating container: %s", err)
	}

	if cont.GetElem() != nil {
		t.Error("Non nil return on empty elems")
	}

	for i := 0; i < elems_len; i++ {

		var testing_str Link = Link(fmt.Sprintf("hi test me %d", i))

		cont.Add(testing_str)

	}

	counter := make(map[string]int)

	// вызвали дохера раз
	for i := 0; i < 10000; i++ {
		el := cont.GetElem()
		if el == nil {
			t.Fatal("Nil pointer to elem!")
		}
		counter[string(el.Item().(Link))]++
	}

	if len(counter) != elems_len {
		t.Fatalf("Invalid get elems, expect %d got %d, container len : %d", elems_len, len(counter), cont.Len())
	}

	f := true
	prev := 0
	// проверяем что распределение равномерное
	for k, val := range counter {

		if f {
			prev = val
			f = false
		} else {
			if prev != val {
				t.Fatalf("Invalid balance: %q %d <> %d", k, prev, val)
			}
		}

	}

	t.Logf("All items: %v", counter)
	t.Logf("All items has weight: %d", prev)

}
