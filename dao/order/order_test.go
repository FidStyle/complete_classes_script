package order

import (
	testarg "compete_classes_script/pkg/test_arg"
	"fmt"
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {
	tx := testarg.Tx()

	tx.AutoMigrate(&Order{})

	od, err := CreateOrder(tx, &Order{
		Pw:        "123456",
		Account:   "654321",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	defer tx.Delete(od)

	var findOd Order
	if err := tx.Where("id = ?", od.ID).Find(&findOd).Error; err != nil {
		t.Error(err)
	}

	if findOd.ID != od.ID {
		t.Errorf("id: %v:%v", od.ID, findOd.ID)
	}
	if findOd.Pw != "123456" {
		t.Errorf("pw: %v:%v", "123465", findOd.Pw)
	}
	if findOd.Account != "654321" {
		t.Errorf("account: %v:%v", "654321", findOd.Account)
	}
}

func TestUpdateOrderScoreByID(t *testing.T) {
	tx := testarg.Tx()
	order, err := CreateOrder(tx, &Order{
		Pw:        "123456",
		Account:   "654321",
		CreatedAt: time.Now(),
		F:         88,
	})
	if err != nil {
		t.Error(err)
	}

	if _, err := UpdateOrderScoreByID(tx, order.ID, "f", 50); err != nil {
		t.Error(err)
	}

	order2, err := GetOrderByID(tx, order.ID)
	if err != nil {
		t.Error(err)
	}
	if len(order2) != 1 {
		t.Errorf("len(order2): %v:%v", 1, len(order2))
	}

	if order2[0].F != 88-50 {
		t.Errorf("F: %v:%v", 88-50, order2[0].F)
	}

	if _, err := UpdateOrderScoreByID(tx, order.ID, "f", 39); err != nil {
		t.Error(err)
	}
	order3, err := GetOrderByID(tx, order.ID)
	if err != nil {
		t.Error(err)
	}
	if len(order3) != 1 {
		t.Errorf("len(order3): %v:%v", 1, len(order3))
	}

	if order3[0].F != 0 {
		t.Errorf("F: %v:%v", 0, order3[0].F)
	}

}

func TestUpdateOrderSchedulingByID(t *testing.T) {
	tx := testarg.Tx()
	order, err := CreateOrder(tx, &Order{
		Pw:        "123456",
		Account:   "654321",
		CreatedAt: time.Now(),
		F:         88,
		Condition: "wait",
	})
	if err != nil {
		t.Error(err)
	}

	if _, err := UpdateOrderConditionByID(tx, order.ID, ConditionScheduling); err != nil {
		t.Error(err)
	}

	order2, err := GetOrderByID(tx, order.ID)
	if err != nil {
		t.Error(err)
	}
	if len(order2) != 1 {
		t.Errorf("len(order2): %v:%v", 1, len(order2))
	}

	if order2[0].Condition != ConditionScheduling {
		t.Errorf("Condition: %v:%v", ConditionScheduling, order2[0].Condition)
	}
}

func TestUpdateOrderSuccessAtByID(t *testing.T) {
	tx := testarg.Tx()
	order, err := CreateOrder(tx, &Order{
		Pw:        "123456",
		Account:   "654321",
		CreatedAt: time.Now(),
		F:         88,
		Condition: ConditionWait,
	})
	if err != nil {
		t.Error(err)
	}

	if _, err := UpdateOrderSuccessAtByID(tx, order.ID, time.Now()); err != nil {
		t.Error(err)
	}

	order2, err := GetOrderByID(tx, order.ID)
	if err != nil {
		t.Error(err)
	}
	if len(order2) != 1 {
		t.Errorf("len(order2): %v:%v", 1, len(order2))
	}

	if order2[0].SuccessAt.Equal(time.Time{}) {
		t.Error("SuccessAt get empty time")
	}
}

func TestUpdateOrderSpecifyPublicByID(t *testing.T) {
	tx := testarg.Tx()

	order, err := CreateOrder(tx, &Order{
		SpecifyPublic: "日语初级 英语口语 武经七书",
	})
	if err != nil {
		t.Error(err)
	}

	order2, err := UpdateOrderSpecifyPublicByID(tx, order.ID, "英语口语")
	if err != nil {
		t.Error(err)
	}

	if order2.SpecifyPublic != fmt.Sprintf("日语初级 武经七书") {
		t.Errorf("order2.SpecifyPublic: %v:%v", fmt.Sprintln("日语初级 武经七书"), order2.SpecifyPublic)
	}
}

func TestGetOrderByCreater(t *testing.T) {
	tx := testarg.Tx()

	orders := []Order{
		{Condition: ConditionScheduling}, {Condition: ConditionError}, {Condition: ConditionSuccess}, {Condition: ConditionScheduling},
		{Condition: ConditionWait}, {Condition: ConditionSlow}, {Condition: ConditionWait, SuccessAt: time.Now()},
	}

	for i := 0; i < len(orders); i++ {
		createOrder := orders[i]
		createOrder.Creater = "test"
		order, err := CreateOrder(tx, &createOrder)
		if err != nil {
			t.Error(err)
		}
		if _, err := UpdateOrderConditionByID(tx, order.ID, orders[i].Condition); err != nil {
			t.Error(err)
		}
	}

	res, err := GetOrderByCreater(tx, 3, 1, "test", false, "false")
	if err != nil {
		t.Error(err)
	}
	if len(res) != 3 {
		t.Errorf("len(res): %v:%v", 3, len(res))
	}

	if res[0].Condition != ConditionScheduling {
		t.Errorf("res[0].Condition: %v:%v", ConditionScheduling, res[0].Condition)
	}
	if res[1].Condition != ConditionScheduling {
		t.Errorf("res[1].Condition: %v:%v", ConditionScheduling, res[1].Condition)
	}
	if res[2].Condition != ConditionSlow {
		t.Errorf("res[2].Condition: %v:%v", ConditionSlow, res[2].Condition)
	}

	// res, err = GetOrderByCreater(tx, 10, 0, "test", true)
	// if err != nil {
	// 	t.Error(err)
	// }
	// if len(res) != 1 {
	// 	t.Errorf("len(res): %v:%v", 1, len(res))
	// }

	// fmt.Println(res[0].CreatedAt)
}

func TestUpdateOrderInfoByID(t *testing.T) {
	tx := testarg.Tx()

	order, err := CreateOrder(tx, &Order{
		Info: true,
	})
	if err != nil {
		t.Error(err)
	}

	newOrder, err := UpdateOrderInfoByID(tx, order.ID, false)
	if err != nil {
		t.Error(err)
	}
	if newOrder.Info {
		t.Errorf("info: %v:%v", false, newOrder.Info)
	}
}

func TestDeleteOrderByID(t *testing.T) {
	tx := testarg.Tx()

	order, err := CreateOrder(tx, &Order{
		Info: true,
	})
	if err != nil {
		t.Error(err)
	}

	if err := DeleteOrderByID(tx, order.ID); err != nil {
		t.Error(err)
	}

	res, _ := GetOrderByID(tx, order.ID)
	// if err != nil {
	// 	t.Error(err)
	// }

	if len(res) != 0 {
		t.Errorf("len(res): %v:%v", 0, len(res))
	}
}
