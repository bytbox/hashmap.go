package hashmap

import (
	"testing"
)

func TestSimple(t *testing.T) {
	m := New()
	if _, b := m.Get(1); b {
		t.Errorf("Get(1) returned true; false expected")
	}
	m.Set(2, "hi")
	if _, b := m.Get(1); b {
		t.Errorf("Get(1) returned true; false expected")
	}
	v, b := m.Get(2)
	if !b {
		t.Errorf("Get(2) returned false; true expected")
	}
	if v != "hi" {
		t.Errorf("Get(2) returned %#v; %#v expected", v, "hi")
	}
	m.Del(2)
	if _, b := m.Get(2); b {
		t.Errorf("Get(2) returned true after Del(2); false expected")
	}
}

func TestCap(t *testing.T) {
	for i := 5; i < 20; i++ {
		if x := NewCap(i).Cap(); x != i {
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

func BenchmarkRawGetFail(b *testing.B) {
	m := newRaw()
	for i := 0; i < b.N; i++ {
		_ = m[5]
	}
}

func BenchmarkGetFail(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		_, _ = m.Get(5)
	}
}

func BenchmarkRawSet(b *testing.B) {
	m := newRaw()
	for i := 0; i < b.N; i++ {
		m[0] = i
	}
}

func BenchmarkRawSetIncremental(b *testing.B) {
	m := newRaw()
	for i := 0; i < b.N; i++ {
		m[i] = i
	}
}
