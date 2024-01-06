package service

import (
	"compete_classes_script/dao/user"
	baseresp "compete_classes_script/pkg/base_resp"
	"compete_classes_script/service/types"
	"context"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserServer struct {
	ctx context.Context
	tx  *gorm.DB
	rtx *redis.Client
}

func NewUserServer(ctx context.Context, tx *gorm.DB, rtx *redis.Client) *UserServer {
	return &UserServer{
		ctx: ctx,
		tx:  tx,
		rtx: rtx,
	}
}

func (s *UserServer) Login(any *types.LoginReq) *types.LoginResp {
	users, err := user.GetUserByAccount(s.tx, any.Account)
	if err != nil {
		return &types.LoginResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}
	if len(users) == 0 {
		return &types.LoginResp{
			BaseResp: *baseresp.ErrorResp(baseresp.ErrLoginUnknownAccount),
		}
	}

	if users[0].Pw != any.Pw {
		return &types.LoginResp{
			BaseResp: *baseresp.ErrorResp(baseresp.ErrLoginWrongPw),
		}
	}

	token, err := user.CreateAccountToken(s.rtx, uuid.New().String(), any.Account)
	if err != nil {
		return &types.LoginResp{
			BaseResp: *baseresp.ErrorResp(err),
		}
	}

	return &types.LoginResp{
		BaseResp: *baseresp.SuccessResp(),
		Token:    token,
	}
}
