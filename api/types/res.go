package types

import (
	"compete_classes_script/dao/order"
	baseresp "compete_classes_script/pkg/base_resp"
)

type CreateOrderResp struct {
	baseresp.BaseResp
}

type GetOrderByCreaterResp struct {
	baseresp.BaseResp
	Orders []*order.Order `json:"orders"`
}

type LoginResp struct {
	baseresp.BaseResp
	Token string `json:"token"`
}
