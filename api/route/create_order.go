package route

import (
	"compete_classes_script/api/svc"
	"compete_classes_script/api/types"
	baseresp "compete_classes_script/pkg/base_resp"
	"compete_classes_script/pkg/logger"
	castas "compete_classes_script/pkg/utils/cast/cast_as"
	"compete_classes_script/service"
	"context"

	"github.com/gin-gonic/gin"
)

func CreateOrder(svctx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.CreateOrderReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(200, baseresp.ErrorResp(baseresp.ErrInvalidArgument))
			return
		}

		if err := service.NewOrderServer(context.Background(), svctx.Tx, svctx.Rtx).CreateOrder(castas.CastCreateOrderReq(&req)); err != nil {
			logger.Error(err)
			c.JSON(200, baseresp.ErrorResp(err))
			return
		}

		c.JSON(200, baseresp.SuccessResp())
	}
}
