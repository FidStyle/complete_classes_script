package route

import (
	"compete_classes_script/api/middleware"
	"compete_classes_script/api/svc"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, svctx *svc.ServiceContext) {
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})
	r.POST("/api/order", middleware.TokenVerify(svctx.Rtx), CreateOrder(svctx))
	r.POST("/api/login", Login(svctx))
	r.POST("/api/get/order", middleware.TokenVerify(svctx.Rtx), GetOrderByCreater(svctx))
}
