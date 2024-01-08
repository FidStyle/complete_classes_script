package order

import (
	baseresp "compete_classes_script/pkg/base_resp"
	"compete_classes_script/pkg/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	ConditionWait       = "wait"
	ConditionScheduling = "scheduling"
	ConditionError      = "error"
	ConditionSlow       = "slow"
	ConditionSuccess    = "success"
)

func CreateOrder(tx *gorm.DB, model *Order) (*Order, error) {
	model.CreatedAt = time.Now()
	model.Condition = ConditionWait
	model.Info = false
	if err := tx.Model(&Order{}).Omit("success_at").Create(&model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func GetOrderByID(tx *gorm.DB, id int) ([]*Order, error) {
	res := []*Order{}
	if err := tx.Model(&Order{}).Where("id = ?", id).First(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func GetOrderByKind(tx *gorm.DB, kind string) ([]*Order, error) {
	res := []*Order{}
	if err := tx.Model(&Order{}).Order("created_at").Where(fmt.Sprintf("%v != 0", kind)).First(res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func GetTheOldestOrder(tx *gorm.DB, limit int) ([]*Order, error) {
	res := []*Order{}
	if err := tx.Model(&Order{}).Order("created_at").Where("success_at is null").Where("condition = ?", ConditionWait).Limit(limit).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func GetOrderByCreater(tx *gorm.DB, limit, offset int, creater string, finish bool, info string) ([]*Order, error) {
	if finish {
		return getOrderByCreaterFinish(tx, limit, offset, creater, info)
	} else {
		return getOrderByCreaterUnfinish(tx, limit, offset, creater)
	}
}

func getOrderByCreaterUnfinish(tx *gorm.DB, limit int, offset int, creater string) ([]*Order, error) {
	res := []*Order{}
	sql := "select * from orders where creater = ? and success_at is null order by case when condition = 'error' then 1 when condition = 'scheduling' then 2 when condition = 'slow' then 3 else 4 end, created_at DESC limit ? offset ?"

	if err := tx.Raw(sql, creater, limit, offset).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func getOrderByCreaterFinish(tx *gorm.DB, limit int, offset int, creater string, info string) ([]*Order, error) {
	res := []*Order{}
	sql := "select * from orders where creater = ? and success_at is not null and info = ? order by success_at DESC limit ? offset ?"
	if info == "ignore" {
		sql = "select * from orders where creater = ? and success_at is not null order by success_at DESC limit ? offset ?"
		if err := tx.Raw(sql, creater, limit, offset).Find(&res).Error; err != nil {
			return nil, err
		}
		return res, nil
	}

	arg := false
	if info == "true" {
		arg = true
	}
	if err := tx.Raw(sql, creater, arg, limit, offset).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateOrderConditionByID(tx *gorm.DB, id int, val string) (*Order, error) {
	if err := tx.Model(&Order{}).Where("id = ?", id).Update("condition", val).Error; err != nil {
		return nil, err
	}

	res := &Order{}
	if err := tx.Model(&Order{}).Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateOrderSuccessAtByID(tx *gorm.DB, id int, t time.Time) (*Order, error) {
	if err := tx.Model(&Order{}).Where("id = ?", id).Update("success_at", t).Error; err != nil {
		return nil, err
	}

	res := &Order{}
	if err := tx.Model(&Order{}).Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateOrderScoreByID(tx *gorm.DB, id int, kind string, decrease float64) (*Order, error) {
	res := &Order{}
	if err := tx.Model(&Order{}).Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}

	score := res.GetScoreByKind(kind)
	score -= decrease
	if score < 0 {
		score = 0
	}

	if err := tx.Model(&Order{}).Where("id = ?", id).Update(kind, score).Error; err != nil {
		return nil, err
	}

	if err := tx.Model(&Order{}).Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateAllOrderToNotScheduling(tx *gorm.DB) error {
	if err := tx.Model(&Order{}).Where("condition = ?", ConditionScheduling).Or("condition = ?", ConditionSlow).Update("condition", ConditionWait).Error; err != nil {
		return err
	}

	return nil
}

func UpdateOrderSpecifyPublicByID(tx *gorm.DB, id int, name string) (*Order, error) {
	res := &Order{}
	if err := tx.Model(&Order{}).Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}

	classes := utils.CastStringToSlice(res.SpecifyPublic)
	find := false
	for i := 0; i < len(classes); i++ {
		if classes[i] == name {
			classes = append(classes[:i], classes[i+1:]...)
			find = true
			break
		}
	}

	if !find {
		return nil, baseresp.ErrSpecifyPublicNotFound
	}

	newClasses := utils.CastSliceToString(classes)

	if err := tx.Model(&Order{}).Where("id = ?", id).Update("specify_public", newClasses).Error; err != nil {
		return nil, err
	}

	if err := tx.Model(&Order{}).Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateOrderSpecifyProfessionalByID(tx *gorm.DB, id int, name string) (*Order, error) {
	res := &Order{}
	if err := tx.Model(&Order{}).Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}

	classes := utils.CastStringToSlice(res.SpecifyProfessional)
	find := false
	for i := 0; i < len(classes); i++ {
		if classes[i] == name {
			classes = append(classes[:i], classes[i+1:]...)
			find = true
			break
		}
	}

	if !find {
		return nil, baseresp.ErrSpecifyPublicNotFound
	}

	newClasses := utils.CastSliceToString(classes)

	if err := tx.Model(&Order{}).Where("id = ?", id).Update("specify_professional", newClasses).Error; err != nil {
		return nil, err
	}

	if err := tx.Model(&Order{}).Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateOrderInfoByID(tx *gorm.DB, id int, info bool) (*Order, error) {
	if err := tx.Model(&Order{}).Where("id = ?", id).Update("info", info).Error; err != nil {
		return nil, err
	}

	res := &Order{}
	if err := tx.Model(&Order{}).Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func DeleteOrderByID(tx *gorm.DB, id int) error {
	if err := tx.Model(&Order{}).Where("id = ?", id).Delete(&Order{}).Error; err != nil {
		return err
	}

	return nil
}
