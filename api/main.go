package main

import (
	"compete_classes_script/api/config"
	"compete_classes_script/api/route"
	"compete_classes_script/api/svc"
	"compete_classes_script/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	c := config.NewConfig("../../config.yaml")
	logger.Init(c.AccessLog, c.ErrorLog)
	svctx := svc.NewServiceContext(c)

	// competer := competer.NewCompeter(svctx)
	// competer.Start()

	route.RegisterRoutes(r, svctx)

	r.Run(":8000")
}
