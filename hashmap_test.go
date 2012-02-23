package hashmap

import (
	"testing"
)

// We do not attempt to test for race conditions. A pseudo-formal proof of
// correctness is provided in the comments to waitmap.go - that will have to
// do. These tests only test single-threaded functionality.
//
// In the future, it might be wise to create a script to insert calls to
// runtime.Gosched() between every pair of lines, which should (coupled with
// GOMAXPROCS>100) be sufficient to catch most problems.

func TestSimple(t *testing.T) {
	m := New()
	if m.Get(1) != nil {
		t.Errorf("Get(1) unexpectedly returned non-nil")
	}
	m.Set(2, "hi")
	if m.Get(1) != nil {
		t.Errorf("Get(1) unexpectedly returned non-nil")
	}
	v := m.Get(2)
	if v != "hi" {
		t.Errorf("Get(2) returned %#v; %#v expected", v, "hi")
	}
	m.Del(2)
	if m.Get(2) != nil {
		t.Errorf("Get(2) unexpectedly returned non-nil")
	}
}

func TestCapacity(t *testing.T) {
	for i := uint32(5); i < uint32(20); i++ {
		if x := NewCap(i).Capacity(); x != i {
			t.Errorf("NewCap(%d).Cap() returned %d", x)
		}
	}
}

func newRaw() map[interface{}]interface{} {
	return make(map[interface{}]interface{}, defCap)
}

func BenchmarkRawCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = newRaw()
	}
}

func BenchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

func BenchmarkRawCreateLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = make(map[interface{}]interface{}, 5000)
	}
}

func BenchmarkCreateLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewCap(5000)
	}
}

func BenchmarkRawLen(b *testing.B) {
	m := newRaw()
	for i := 0; i < b.N; i++ {
		_ = len(m)
	}
}

func BenchmarkLen(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		_ = m.Size()
	}
}

func BenchmarkRawGetFail(b *testing.B) {
	m := newRaw()
	for i := 0; i < b.N; i++ {
		_ = m[5]
	}
}

func BenchmarkGetFail(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		_ = m.Get(5)
	}
}

func BenchmarkRawSet(b *testing.B) {
	m := newRaw()
	for i := 0; i < b.N; i++ {
		m[0] = i
	}
}

func BenchmarkSet(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Set(0, i)
	}
}

func BenchmarkRawSetIncremental(b *testing.B) {
	m := newRaw()
	for i := 0; i < b.N; i++ {
		m[i] = i
	}
}

func BenchmarkSetIncremental(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}

