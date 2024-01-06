package service

import (
	"compete_classes_script/dao/order"
	"compete_classes_script/dao/user"
	castsd "compete_classes_script/pkg/utils/cast/cast_sd"
	"compete_classes_script/service/types"
	"context"

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
