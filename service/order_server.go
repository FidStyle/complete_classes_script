package service

import (
	"compete_classes_script/dao/order"
	"compete_classes_script/dao/user"
	baseresp "compete_classes_script/pkg/base_resp"
	castsd "compete_classes_script/pkg/utils/cast/cast_sd"
	"compete_classes_script/service/types"
	"context"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type OrderServer struct {
	ctx context.Context
	tx  *gorm.DB
	rtx *redis.Client
}

func NewOrderServer(ctx context.Context, tx *gorm.DB, rtx *redis.Client) *OrderServer {
	return &OrderServer{
		ctx: ctx,
		tx:  tx,
		rtx: rtx,
	}
}

func (s *OrderServer) CreateOrder(any *types.CreateOrderReq) error {
	account, err := user.GetAccountByToken(s.rtx, any.Token)
	if err != nil {
		return err
	}

	req := castsd.CastCreateOrderReqToOrder(any)
	req.Creater = account
	if _, err := order.CreateOrder(s.tx, req); err != nil {
		return err
	}

	return nil
}

func (s *OrderServer) GetOrderByCreater(any *types.GetOrderByCreaterReq) *types.GetOrderByCreaterResp {
	account, err := user.GetAccountByToken(s.rtx, any.Token)
	if err != nil {
		return &types.GetOrderByCreaterResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}
	res, err := order.GetOrderByCreater(s.tx, any.Limit, any.Offset, account, any.Finish, any.Info)
	if err != nil {
		return &types.GetOrderByCreaterResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}

	for i := 0; i < len(res); i++ {
		if !res[i].SuccessAt.Equal(time.Time{}) {
			res[i].SuccessAt = res[i].SuccessAt.Add(8 * time.Hour)
		}
		if !res[i].CreatedAt.Equal(time.Time{}) {
			res[i].CreatedAt = res[i].CreatedAt.Add(8 * time.Hour)
		}
	}

	return &types.GetOrderByCreaterResp{
		BaseResp: *baseresp.SuccessResp(),
		Orders:   res,
	}
}

func (s *OrderServer) UpdateOrderInfoByID(any *types.UpdateOrderInfoByIDReq) *types.UpdateOrderInfoByIDResp {
	account, err := user.GetAccountByToken(s.rtx, any.Token)
	if err != nil {
		return &types.UpdateOrderInfoByIDResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}

	orders, err := order.GetOrderByID(s.tx, any.ID)
	if err != nil {
		return &types.UpdateOrderInfoByIDResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}
	if len(orders) == 0 {
		return &types.UpdateOrderInfoByIDResp{
			BaseResp: *baseresp.ErrorResp(baseresp.ErrOrderIDNotExist),
		}
	}

	if orders[0].Creater != account {
		return &types.UpdateOrderInfoByIDResp{
			BaseResp: *baseresp.ErrorResp(baseresp.ErrAuthInvalid),
		}
	}

	if orders[0].SuccessAt.Equal(time.Time{}) {
		return &types.UpdateOrderInfoByIDResp{
			BaseResp: *baseresp.ErrorResp(baseresp.ErrInfoUnsuccessOrder),
		}
	}

	if _, err := order.UpdateOrderInfoByID(s.tx, any.ID, any.Info); err != nil {
		return &types.UpdateOrderInfoByIDResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}

	return &types.UpdateOrderInfoByIDResp{
		BaseResp: *baseresp.SuccessResp(),
	}
}

func (s *OrderServer) DeleteOrderByID(any *types.DeleteOrderByIDReq) *types.DeleteOrderByIDResp {
	account, err := user.GetAccountByToken(s.rtx, any.Token)
	if err != nil {
		return &types.DeleteOrderByIDResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}

	orders, err := order.GetOrderByID(s.tx, any.ID)
	if err != nil {
		return &types.DeleteOrderByIDResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}
	if len(orders) == 0 {
		return &types.DeleteOrderByIDResp{
			BaseResp: *baseresp.ErrorResp(baseresp.ErrOrderIDNotExist),
		}
	}

	if orders[0].Creater != account {
		return &types.DeleteOrderByIDResp{
			BaseResp: *baseresp.ErrorResp(baseresp.ErrAuthInvalid),
		}
	}

	if err := order.DeleteOrderByID(s.tx, any.ID); err != nil {
		return &types.DeleteOrderByIDResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}

	return &types.DeleteOrderByIDResp{
		BaseResp: *baseresp.SuccessResp(),
	}
}
