package reset

import "testing"

type TestStruct struct {
	Value int
}

func (ts *TestStruct) Reset() {
	ts.Value = 0
}
func NewTestStruct() *TestStruct {
	return &TestStruct{Value: 42}
}

func TestPoolGet(t *testing.T) {

	pool := NewPool(NewTestStruct)
	obj := pool.Get()

	if obj == nil {
		t.Error("Expected non-nil object from pool")
		return
	}

	if obj.Value != 42 {
		t.Errorf("Expected initial value 42, got %d", obj.Value)
	}
}

func TestPoolPut(t *testing.T) {
	// Создаем пул
	pool := NewPool(NewTestStruct)
	obj := pool.Get()

	pool.Put(obj)
	if obj == nil {
		t.Error("Expected non-nil object from pool")
		return
	}

	if obj.Value != 0 {
		t.Errorf("Expected initial value 42, got %d", obj.Value)
	}
}
