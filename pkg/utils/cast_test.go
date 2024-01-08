package utils

import "testing"

func TestCastStringToSlice(t *testing.T) {
	res := CastStringToSlice("A B")
	if len(res) != 2 {
		t.Errorf("len(res): %v:%v", 2, len(res))
	}

	if res[0] != "A" || res[1] != "B" {
		t.Errorf("res: %v:%v", []string{"A", "B"}, res)
	}
}
