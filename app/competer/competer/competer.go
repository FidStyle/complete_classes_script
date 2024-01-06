package competer

import (
	"compete_classes_script/api/svc"
	"compete_classes_script/dao/order"
	baseresp "compete_classes_script/pkg/base_resp"
	"compete_classes_script/pkg/heu"
	"compete_classes_script/pkg/logger"
	"compete_classes_script/pkg/utils"
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Competer struct {
	svctx *svc.ServiceContext

	disPatcher *DisPatcher
	taskChan   chan *task
}

func NewCompeter(svctx *svc.ServiceContext) *Competer {
	taskChan := make(chan *task)
	return &Competer{
		svctx:      svctx,
		taskChan:   taskChan,
		disPatcher: NewDisPatcher(svctx, taskChan),
	}
}

func (c *Competer) Start() {
	logger.Info("Competer Start ...")
	if err := c.CleanData(); err != nil {
		panic(err)
	}
	go c.disPatcher.Start()
	go c.SolveTask()
}

func (c *Competer) CleanData() error {
	return order.UpdateAllOrderToNotScheduling(c.svctx.Tx)
}

func (c *Competer) SolveTask() {
	for t := range c.taskChan {
		go c.solveTask(t, c.svctx.Cfg.CaptchaToken)
	}
}

func (c *Competer) solveTask(t *task, token string) {
	c.disPatcher.meta.UpdateCount(1)
	defer func() {
		c.disPatcher.meta.UpdateCount(-1)
		if _, err := order.UpdateOrderSuccessAtByID(c.svctx.Tx, t.order.ID, time.Now()); err != nil {
			logger.Error(err)
		}
		if _, err := order.UpdateOrderConditionByID(c.svctx.Tx, t.order.ID, order.ConditionWait); err != nil {
			logger.Error(err)
		}
	}()

	xfs := t.order.FindUnZeroScore()
	wg := sync.WaitGroup{}
	if _, err := order.UpdateOrderConditionByID(c.svctx.Tx, t.order.ID, order.ConditionScheduling); err != nil {
		logger.Error(err)
	}

	logger.Infof("%v start", t.order.Account)
	for _, xf := range xfs {
		wg.Add(1)
		if strings.Contains(xf, "_n") {
			c.solveTaskByKindUntilSuccess(t, xf, &wg, NetworkWrapFunc)
		} else if strings.Contains(xf, "SpecifyPublic") {
			classes := utils.CastStringToSlice(t.order.SpecifyPublic)
			c.solveTaskBySpecifyUntilSuccess(t, classes, &wg)
		} else if strings.Contains(xf, "SpecifyProfessional") {
			classes := utils.CastStringToSlice(t.order.SpecifyProfessional)
			c.solveTaskBySpecifyProfessionalUntilSuccess(t, classes, &wg)
		} else {
			c.solveTaskByKindUntilSuccess(t, xf, &wg)
		}
	}

	wg.Wait()

	logger.Infof("[final] %v success with harvesting %v", t.order.Account, t.harvest)
}

func (c *Competer) solveTaskByKindUntilSuccess(t *task, kind string, wg *sync.WaitGroup, wraps ...func(string) bool) {
	solve := func() error {
		token, batchID := heu.LoginUntilSuccess(t.order.Account, t.order.Pw, c.svctx.Cfg.CaptchaToken)
		for {
			ok, err := c.solveTaskByKind(t, kind, token, batchID, wraps...)
			if err != nil {
				return err
			}

			time.Sleep(time.Duration(c.svctx.Cfg.EverySolveSleepTime) * time.Millisecond)

			if ok {
				logger.Infof("%v with kind %v successes and it harvests %v", t.order.Account, kind, t.harvest)
				return nil
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go CountTime(ctx, t, c.svctx.Tx)

	defer wg.Done()
	for {
		err := solve()
		if err != nil {
			if errors.Is(err, baseresp.ErrHeuReLogin) {
				continue
			} else if errors.Is(err, baseresp.ErrHeuFullClasses) {
				logger.Infof("[InfoError] %v is full but need %v score of %v kind with harvesting %v", t.order.Account, t.order.GetScoreByKind(kind), kind, t.harvest)
				return
			} else if errors.Is(err, baseresp.ErrKindNotFound) || errors.Is(err, baseresp.ErrKindNotFoundInWrap) {
				logger.Infof("[InfoError] %v expect error: %v", err, t.order.Account)
				if _, err := order.UpdateOrderConditionByID(c.svctx.Tx, t.order.ID, order.ConditionError); err != nil {
					logger.Error(err)
				}
				return
			}
			logger.Errorf("solveTaskByKind failed because of: %v", err.Error())
			time.Sleep(time.Duration(c.svctx.Cfg.ErrorSleepTime) * time.Second)
		} else {
			return
		}
	}
}

func (c *Competer) solveTaskByKind(t *task, kind, token, batchID string, wraps ...func(string) bool) (bool, error) {
	rows, err := heu.GetList(utils.ChangeNToNormal(kind), token, batchID)
	if err != nil {
		return false, err
	}
	if len(rows.Data.Rows) == 0 {
		return false, baseresp.ErrKindNotFound
	}

	find := false
	for _, row := range rows.Data.Rows {
		if !WrapCheck(row.KCM, wraps...) {
			continue
		}
		find = true
		if row.KRL > row.YXRS {
			ok, err := heu.AddClassVolunteer(row.JXBID, row.SecretVal, token, batchID)
			if err != nil {
				return false, err
			}

			if ok {
				xf, err := strconv.ParseFloat(row.XF, 64)
				if err != nil {
					return false, err
				}

				t.harvest = append(t.harvest, row.KCM)
				updateOrder, err := order.UpdateOrderScoreByID(c.svctx.Tx, t.order.ID, kind, xf)
				if err != nil {
					return false, err
				}

				logger.Infof("%v harvest %v with %v score, and it still needs to harvest %v more scores",
					t.order.Account, row.KCM, row.XF, updateOrder.GetScoreByKind(kind))

				t.order = updateOrder

				if updateOrder.GetScoreByKind(kind) == 0 {
					return true, nil
				}
			}
		}
	}

	if !find {
		return false, baseresp.ErrKindNotFoundInWrap
	}

	return false, nil
}

func (c *Competer) solveTaskBySpecifyUntilSuccess(t *task, fuzzyNames []string, wg *sync.WaitGroup) {
	solve := func(fuzzyName string) error {
		token, batchID := heu.LoginUntilSuccess(t.order.Account, t.order.Pw, c.svctx.Cfg.CaptchaToken)
		for {
			ok, err := c.solveTaskBySpecify(t, fuzzyName, token, batchID)
			if err != nil {
				return err
			}

			time.Sleep(time.Duration(c.svctx.Cfg.EverySolveSleepTime) * time.Millisecond)

			if ok {
				logger.Infof("%v of specifyPublic successes and it harvests %v", t.order.Account, t.harvest)
				return nil
			}
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go CountTime(ctx, t, c.svctx.Tx)

	defer wg.Done()
	for _, fuzzyName := range fuzzyNames {
		for {
			err := solve(fuzzyName)
			if err != nil {
				if errors.Is(err, baseresp.ErrHeuReLogin) {
					continue
				} else if errors.Is(err, baseresp.ErrHeuFullClasses) {
					logger.Infof("[InfoError] %v is full but need %v unfinished with harvesting %v", t.order.Account, t.order.SpecifyPublic, t.harvest)
					return
				} else if errors.Is(err, baseresp.ErrKindNotFound) || errors.Is(err, baseresp.ErrKindNotFoundInWrap) {
					logger.Infof("[InfoError] %v expect error: %v", err, t.order.Account)
					if _, err := order.UpdateOrderConditionByID(c.svctx.Tx, t.order.ID, order.ConditionError); err != nil {
						logger.Error(err)
					}
					return
				}
				logger.Errorf("solveTaskByKind failed because of: %v", err.Error())
				time.Sleep(time.Duration(c.svctx.Cfg.ErrorSleepTime) * time.Second)
			} else {
				break
			}
		}
	}
}

func (c *Competer) solveTaskBySpecify(t *task, fuzzyName string, token, batchID string) (bool, error) {
	rows, err := heu.GetListByFuzzyName(fuzzyName, token, batchID)
	if err != nil {
		return false, err
	}
	if len(rows.Data.Rows) == 0 {
		return false, baseresp.ErrFuzzyNameNotFound
	}

	for _, row := range rows.Data.Rows {
		ok, err := heu.AddClassVolunteer(row.JXBID, row.SecretVal, token, batchID)
		if err != nil {
			return false, err
		}

		if ok {
			t.harvest = append(t.harvest, row.KCM)
			updateOrder, err := order.UpdateOrderSpecifyPublicByID(c.svctx.Tx, t.order.ID, fuzzyName)
			if err != nil {
				return false, err
			}

			logger.Infof("%v harvest %v and it still needs to harvest %v", t.order.Account, row.KCM, updateOrder.SpecifyPublic)

			t.order = updateOrder
			return true, nil
		}
	}
	return false, nil
}

func (c *Competer) solveTaskBySpecifyProfessionalUntilSuccess(t *task, fuzzyNames []string, wg *sync.WaitGroup) {
	solve := func(fuzzyName string) error {
		token, batchID := heu.LoginUntilSuccess(t.order.Account, t.order.Pw, c.svctx.Cfg.CaptchaToken)
		for {
			ok, err := c.solveTaskBySpecifyProfessional(t, fuzzyName, token, batchID)
			if err != nil {
				return err
			}

			time.Sleep(time.Duration(c.svctx.Cfg.EverySolveSleepTime) * time.Millisecond)

			if ok {
				logger.Infof("%v of specifyPublicProfessional successes and it harvests %v", t.order.Account, t.harvest)
				return nil
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go CountTime(ctx, t, c.svctx.Tx)

	defer wg.Done()
	for _, fuzzyName := range fuzzyNames {
		for {
			err := solve(fuzzyName)
			if err != nil {
				if errors.Is(err, baseresp.ErrHeuReLogin) {
					continue
				} else if errors.Is(err, baseresp.ErrHeuFullClasses) {
					logger.Infof("[InfoError] %v is full but need %v unfinished with harvesting %v", t.order.Account, t.order.SpecifyPublic, t.harvest)
					return
				} else if errors.Is(err, baseresp.ErrKindNotFound) || errors.Is(err, baseresp.ErrKindNotFoundInWrap) {
					logger.Infof("[InfoError] %v expect error: %v", err, t.order.Account)
					if _, err := order.UpdateOrderConditionByID(c.svctx.Tx, t.order.ID, order.ConditionError); err != nil {
						logger.Error(err)
					}
					return
				}
				logger.Errorf("solveTaskByKind failed because of: %v", err.Error())
				time.Sleep(time.Duration(c.svctx.Cfg.ErrorSleepTime) * time.Second)
			} else {
				break
			}
		}
	}
}

func (c *Competer) solveTaskBySpecifyProfessional(t *task, fuzzyName string, token, batchID string) (bool, error) {
	rows, err := heu.GetListProfessional(token, batchID)
	if err != nil {
		return false, err
	}
	find := false
	for _, row := range rows.Data.Rows {
		if len(row.TcList) == 0 {
			logger.Errorf("solveTaskBySpecifyProfessional len(row.TcList) = 0 account = %v", t.order.Account)
			continue
		}

		for _, tc := range row.TcList {
			if !strings.Contains(tc.KCM, fuzzyName) {
				continue
			}
			find = true
			if tc.KRL > tc.YXRS {
				ok, err := heu.AddClassVolunteer(tc.JXBID, tc.SecretVal, token, batchID, "professional")
				if err != nil {
					return false, err
				}

				if ok {
					t.harvest = append(t.harvest, tc.KCM)
					updateOrder, err := order.UpdateOrderSpecifyProfessionalByID(c.svctx.Tx, t.order.ID, fuzzyName)
					if err != nil {
						return false, err
					}

					logger.Infof("%v harvest %v and it still needs to harvest %v", t.order.Account, tc.KCM, updateOrder.SpecifyPublic)

					t.order = updateOrder
					return true, nil
				}
			}
		}
	}

	if !find {
		return false, baseresp.ErrFuzzyNameNotFound
	}
	return false, nil
}

func NetworkWrapFunc(KCM string) bool {
	return strings.Contains(KCM, "网络")
}

func WrapCheck(KCM string, wraps ...func(string) bool) bool {
	for _, wrap := range wraps {
		if !wrap(KCM) {
			return false
		}
	}

	return true
}

type task struct {
	order   *order.Order
	harvest []string
}

func NewTask(order *order.Order) *task {
	return &task{
		order:   order,
		harvest: make([]string, 0),
	}
}

func CountTime(ctx context.Context, t *task, tx *gorm.DB) {
	start := time.Now()
	count := 1
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if time.Now().After(start.Add(time.Duration(count) * time.Hour)) {
				logger.Infof("[InfoSlow] %v is slow with %v hours", t.order.Account, count)
				if _, err := order.UpdateOrderConditionByID(tx, t.order.ID, order.ConditionSlow); err != nil {
					logger.Error(err)
				}
				count++
			}
			time.Sleep(10 * time.Minute)
		}
	}
}
