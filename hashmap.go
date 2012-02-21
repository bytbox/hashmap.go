// Package hashmap implements a simple non-blocking (lock-free) hash map.
package hashmap

import (
	r "reflect"
)

var (
	defCap = 32
)

type entry struct {
	k interface{}
	v interface{}
}

type Hashmap struct {
	data []entry
	size int
	capc int
}

func New() *Hashmap {
	return NewCap(defCap)
}

func NewCap(c int) *Hashmap {
	return &Hashmap{
		data: make([]entry, c),
		size: 0,
		capc: c,
	}
}

func (hm *Hashmap) Cap() int { return hm.capc }

func (hm *Hashmap) Size() int { return hm.size }

func (hm *Hashmap) Get(k interface{}) (interface{}, bool) {
	c := hm.Cap()
	h := hash(k, c)
	for i := 0; i < h+c; i++ {
		e := hm.data[i%c]
		if e.k == nil { return nil, false }
		if r.DeepEqual(e.k, k) {
			return e.v, true
		}
	}
	return nil, false
}

func (hm *Hashmap) Set(k interface{}, v interface{}) {
	c := hm.Cap()
	h := hash(k, c)
	for i := 0; i < h+c; i++ {
		e := &hm.data[i%c]
		if e.k == nil {
			e.v = v
			e.k = k
			hm.size++
			return
		}
		if r.DeepEqual(e.k, k) {
			e.v = v
			hm.size++
			return
		}
	}
	panic("Hashmap full")
}

func (hm *Hashmap) Del(k interface{}) bool {
	c := hm.Cap()
	h := hash(k, c)
	for i := 0; i < h+c; i++ {
		e := &hm.data[i%c]
		if r.DeepEqual(e.k, k) {
			e.k = nil
			e.v = nil
			hm.size--
			return true
		}
	}
	return false
}

func (hm *Hashmap) Grow(c int) {

}

func hash(k interface{}, c int) int {
	switch a := k.(type) {
	case int:
		return a % c
	}
	return 0
}
