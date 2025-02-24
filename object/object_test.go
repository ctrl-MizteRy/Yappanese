package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "hello there"}
	hello2 := &String{Value: "hello there"}
	diff1 := &String{Value: "my name is jeff"}
	deff2 := &String{Value: "my name is jeff"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("String with the same content have different hash key")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("String with the same content have different hash key")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("String with different content have the same hash keys")
	}
}
