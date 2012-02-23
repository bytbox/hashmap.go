// Package hashmap implements a simple non-blocking (lock-free) hash map.
//
// The semantics of this hashmap are quite different from the go built-in.
// Aside from thread-safety, any values (including structs and slices) are
// allowed as keys, and a value of nil is indestinguishable from an element not
// being in the map.
package hashmap

import (
	a "sync/atomic"
	r "reflect"
	u "unsafe"
)

var (
	defCap = uint32(32)
)

// A hashmap entry.
type entry struct {
	f uint32
	k interface{}
	v *interface{}
}

type Hashmap struct {
	data []entry
	size uint32
	capc uint32
}

// Creates a hashmap with the default initial capacity (32). The hashmap will
// be grown to twice its size when the number of entries reaches 50% of the
// capacity.
func New() *Hashmap {
	return NewCap(defCap)
}

// Creates a hashmap with the specified initial capacity.
func NewCap(c uint32) *Hashmap {
	return &Hashmap{
		data: make([]entry, c),
		size: 0,
		capc: uint32(c),
	}
}

// Returns the capacity of this hashmap.
func (hm *Hashmap) Capacity() uint32 {
	return a.LoadUint32(&hm.capc)
}

// Returns the size of this hashmap. This is /not/ the number of elements
// stored - it is the number of keys allocated.
func (hm *Hashmap) Size() uint32 {
	return a.LoadUint32(&hm.size)
}

func (hm *Hashmap) Get(k interface{}) interface{} {
Retry:
	// If the map changes capacity during the operation, we will detect
	// that and retry. Note that the only legal transformation of this
	// value is for it to increase (double), so we can easily check to see
	// if it has changed.
	c := hm.Capacity()

	// Guaranteed to be correct when c is correct.
	h := hash(k, c)

	// Since keys cannot be removed, we know that our target entry lies
	// somewhere 'after' the entry indexed by h, and that no nil keys lie
	// in between.
	for i := uint32(0); i < h+c; i++ {
		// Atomically fetch the data indexed by (i%c). Since there are
		// 16 bytes to fetch, this can't be done in a single atomic
		// operation; however, since a set key can never change, we
		// don't need perfect atomicity when fetching it.

		ef := a.LoadUint32(&hm.data[i%c].f)
		if ef == 0 {
			// The key is not set; there's no such element.
			return nil
		}

		// The key is set and can never change - fetch it
		// non-atomically.
		ek := hm.data[i%c].k

		// Load value atomically
		evp := u.Pointer(hm.data[i%c].v)
		ev := (*interface{})(a.LoadPointer(&evp))

		// Note that at this point, the collect data is incorrect if
		// and only if the capacity has changed (and the map has
		// therefore been expanded). We now check to see if that has
		// happened, and if so, abort.
		if c != hm.Capacity() {
			// This is the part that breaks our per-goroutine
			// progress guarantee.
			goto Retry
		}

		// The data (ek, ev) may now be considered actionable. Since it
		// is stored in local memory, we don't have to worry about
		// making further operations atomic.

		if r.DeepEqual(ek, k) {
			return *ev
		}
	}

	// Really, this should never happen; it just means we've cycled through
	// the whole map, and none of the keys were nil.
	return nil
}

// The key will be (shallowly) copied - the value will not.
func (hm *Hashmap) Set(k interface{}, v interface{}) {
	c := hm.Capacity()
	h := hash(k, c)
	for i := uint32(0); i < h+c; i++ {
		e := &hm.data[i%c]
		if e.k == nil {
			e.f = 1
			e.v = &v
			e.k = k
			hm.size++
			// TODO Grow if necessary
			return
		}
		if r.DeepEqual(e.k, k) {
			e.v = &v
			return
		}
	}
	panic("Hashmap full")
}

// Deletes the given key from the map, and returns true if it existed. Although
// this method is semantically equivalent to m.Set(k, nil), this method will
// not allocate another key if the key is not already in use, and should
// therefore be used instead.
func (hm *Hashmap) Del(k interface{}) bool {
	c := hm.Capacity()
	h := hash(k, c)
	for i := uint32(0); i < h+c; i++ {
		e := &hm.data[i%c]
		if r.DeepEqual(e.k, k) {
			//e.k = nil
			e.f = 0
			e.v = nil
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
