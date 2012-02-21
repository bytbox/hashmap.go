// Package hashmap implements a simple non-blocking (lock-free) hash map.
package hashmap

import (
	a "sync/atomic"
	r "reflect"
)

var (
	defCap = uint32(32)
)

type entry struct {
	k interface{}
	v interface{}
}

type Hashmap struct {
	data []entry
	size uint32
	capc uint32
}

func New() *Hashmap {
	return NewCap(defCap)
}

func NewCap(c uint32) *Hashmap {
	return &Hashmap{
		data: make([]entry, c),
		size: 0,
		capc: uint32(c),
	}
}

// Returns the capacity of this hashmap.
func (hm *Hashmap) Cap() uint32 {
	return a.LoadUint32(&hm.capc)
}

// Returns the size of this hashmap.
func (hm *Hashmap) Size() uint32 {
	return a.LoadUint32(&hm.size)
}

func (hm *Hashmap) Get(k interface{}) (interface{}, bool) {
	c := hm.Cap()
	h := hash(k, c)
	for i := uint32(0); i < h+c; i++ {
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
	for i := uint32(0); i < h+c; i++ {
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
	for i := uint32(0); i < h+c; i++ {
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

func (hm *Hashmap) Grow(c uint32) {
	// TODO
}

func hash(key interface{}, c uint32) uint32 {
	switch k := key.(type) {
	case int:
		return uint32(k) % c
	}
	return uint32(0)
}
