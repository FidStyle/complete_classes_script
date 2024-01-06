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

func Login(svctx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.LoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(200, baseresp.ErrorResp(err))
			return
		}

		resp := service.NewUserServer(context.Background(), svctx.Tx, svctx.Rtx).Login(castas.CastLoginReq(&req))
		c.JSON(200, castas.CastLoginResp(resp))
	}
}
