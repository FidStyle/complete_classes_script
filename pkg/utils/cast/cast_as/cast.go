package castas

import (
	"compete_classes_script/api/types"
	"compete_classes_script/pkg/utils"
	svc "compete_classes_script/service/types"
)

func CastCreateOrderReq(any *types.CreateOrderReq) *svc.CreateOrderReq {
	res := &svc.CreateOrderReq{}

	utils.FillSameField(any, res)
	return res
}

func CastGetOrderByCreaterReq(any *types.GetOrderByCreaterReq) *svc.GetOrderByCreaterReq {
	res := &svc.GetOrderByCreaterReq{}

	utils.FillSameField(any, res)
	return res
}

func CastGetOrderByCreaterResp(any *svc.GetOrderByCreaterResp) *types.GetOrderByCreaterResp {
	res := &types.GetOrderByCreaterResp{}

	utils.FillSameField(any, res)
	return res
}

func CastLoginReq(any *types.LoginReq) *svc.LoginReq {
	res := &svc.LoginReq{}

	utils.FillSameField(any, res)
	return res
}

func CastLoginResp(any *svc.LoginResp) *types.LoginResp {
	res := &types.LoginResp{}

	utils.FillSameField(any, res)
	return res
}
