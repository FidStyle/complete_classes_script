package main

import (
	"compete_classes_script/api/config"
	"compete_classes_script/api/svc"
	"compete_classes_script/app/competer/competer"
	"compete_classes_script/pkg/logger"
)

func main() {
	c := config.NewConfig("../../../config.yaml")
	logger.Init(c.AccessLog, c.ErrorLog)
	svctx := svc.NewServiceContext(c)

	competer := competer.NewCompeter(svctx)
	competer.Start()

	select {}
}
