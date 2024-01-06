package castsd

import (
	"compete_classes_script/dao/order"
	"compete_classes_script/pkg/utils"
	"compete_classes_script/service/types"
)

func CastCreateOrderReqToOrder(any *types.CreateOrderReq) *order.Order {
	res := &order.Order{}

	utils.FillSameField(any, res)

	res.SpecifyProfessional = utils.CastSliceToString(any.SpecifyProfessional)
	res.SpecifyPublic = utils.CastSliceToString(any.SpecifyPublic)

	return res
}
