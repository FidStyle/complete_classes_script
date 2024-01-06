package order

import (
	testarg "compete_classes_script/pkg/test_arg"
	"testing"
)

func TestAutoMig(t *testing.T) {
	tx := testarg.Tx()

	tx.AutoMigrate(&Order{})
}

func TestFindUnZeroScore(t *testing.T) {
	order := &Order{
		ID:   1,
		A:    2,
		F:    3,
		A0_n: 4,
	}

	res := order.FindUnZeroScore()
	if len(res) != 3 {
		t.Errorf("len(res): %v:%v", 3, len(res))
	}

	if res[0] != "a0_n" || res[1] != "a" || res[2] != "f" {
		t.Errorf("res: %v:%v", []string{"a0_n", "a", "f"}, res)
	}
}
