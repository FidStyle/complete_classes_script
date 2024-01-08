package route

import (
	"compete_classes_script/api/svc"
	"compete_classes_script/api/types"
	baseresp "compete_classes_script/pkg/base_resp"
	castas "compete_classes_script/pkg/utils/cast/cast_as"
	"compete_classes_script/service"
	"context"

	"github.com/gin-gonic/gin"
)

func DeleteOrderByID(svctx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.DeleteOrderByIDReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(200, baseresp.ErrorResp(baseresp.ErrInvalidArgument))
			return
		}

		resp := service.NewOrderServer(context.Background(), svctx.Tx, svctx.Rtx).DeleteOrderByID(castas.CastDeleteOrderByIDReq(&req))

		c.JSON(200, resp)
	}
}
