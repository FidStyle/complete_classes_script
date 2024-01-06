package competer

import (
	"compete_classes_script/api/svc"
	"compete_classes_script/dao/order"
	"compete_classes_script/pkg/logger"
	"time"
)

type DisPatcher struct {
	svctx *svc.ServiceContext
	meta  *CompeterMeta

	disPatchTicker time.Ticker
	taskChan       chan *task
}

func NewDisPatcher(svctx *svc.ServiceContext, taskChan chan *task) *DisPatcher {
	return &DisPatcher{
		svctx:    svctx,
		meta:     NewCompeterMeta(),
		taskChan: taskChan,
	}
}

func (c *DisPatcher) Start() {
	go c.disPatch()
}

func (c *DisPatcher) disPatch() {
	c.disPatchTicker = *time.NewTicker(time.Duration(c.svctx.Cfg.DisPatcherTickerTime) * time.Second)

	for range c.disPatchTicker.C {
		count := c.meta.GetCount()
		if count < c.svctx.Cfg.MaxConcurrency {
			orders, err := order.GetTheOldestOrder(c.svctx.Tx, c.svctx.Cfg.MaxConcurrency-count)
			if err != nil {
				logger.Error(err)
				continue
			}

			for _, order := range orders {
				c.taskChan <- NewTask(order)
			}
		}
	}
}
