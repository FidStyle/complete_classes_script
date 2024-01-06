package user

import (
	testarg "compete_classes_script/pkg/test_arg"
	"testing"

	"github.com/google/uuid"
)

func TestCreateAccountToken(t *testing.T) {
	rtx := testarg.Rtx()

	key := uuid.New().String()
	if _, err := CreateAccountToken(rtx, key, "10000000"); err != nil {
		t.Error(err)
	}

	res, err := rtx.Get(key).Result()
	if err != nil {
		t.Error(err)
	}

	if res != "10000000" {
		t.Errorf("res: %v:%v", "10000000", res)
	}

}

func TestGetUserByAccount(t *testing.T) {
	tx := testarg.Tx()

	tx.AutoMigrate(&User{})

	user, err := GetUserByAccount(tx, "10000000")
	if err != nil {
		t.Error(err)
	}
	if len(user) != 1 {
		t.Errorf("len(user): 1:%v", len(user))
	}
	if user[0].Pw != "123456" {
		t.Errorf("Pw: %v:%v", "123456", user[0].Pw)
	}
}
