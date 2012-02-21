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
}

func TestCap(t *testing.T) {
	for i := 5; i < 20; i++ {
		if x := NewCap(i).Cap(); x != i {
			t.Errorf("NewCap(%d).Cap() returned %d", x)
		}
	}
}
