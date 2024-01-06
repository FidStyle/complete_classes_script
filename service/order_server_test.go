package service

import (
	"compete_classes_script/dao/order"
	"compete_classes_script/dao/user"
	testarg "compete_classes_script/pkg/test_arg"
	"compete_classes_script/pkg/utils"
	"compete_classes_script/service/types"
	"context"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	tx := testarg.Tx()
	rtx := testarg.Rtx()
	s := NewOrderServer(context.Background(), tx, rtx)

	_, err := user.CreateAccountToken(rtx, "test in redis", "test")
	if err != nil {
		t.Error(err)
	}

	err = s.CreateOrder(&types.CreateOrderReq{
		Account:             "123456",
		Pw:                  "654321",
		SpecifyPublic:       []string{"孙子兵法"},
		SpecifyProfessional: []string{"常微分", "经济思想史"},
		Token:               "test in redis",
	})
	if err != nil {
		t.Error(err)
	}

	findOrder := &order.Order{}
	if err := tx.Find(findOrder).Error; err != nil {
		t.Error(err)
	}
	defer tx.Delete(findOrder)

	if findOrder.Creater != "test" {
		t.Errorf("Creater: %v:%v", "test", findOrder.Creater)
	}
	if findOrder.Account != "123456" {
		t.Errorf("Account: %v:%v", "123465", findOrder.Account)
	}
	if findOrder.Pw != "654321" {
		t.Errorf("Pw: %v:%v", "654321", findOrder.Pw)
	}
	if findOrder.SpecifyPublic != utils.CastSliceToString([]string{"孙子兵法"}) {
		t.Errorf("SpecifyPublic: %v:%v", utils.CastSliceToString([]string{"孙子兵法"}), findOrder.SpecifyPublic)
	}
	if findOrder.SpecifyProfessional != utils.CastSliceToString([]string{"常微分", "经济思想史"}) {
		t.Errorf("SpecifyProfessional: %v:%v", utils.CastSliceToString([]string{"常微分", "经济思想史"}), findOrder.SpecifyProfessional)
	}
}
