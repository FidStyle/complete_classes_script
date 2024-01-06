package user

import (
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func GetUserByAccount(tx *gorm.DB, account string) ([]*User, error) {
	res := []*User{}
	if err := tx.Model(&User{}).Where("account = ?", account).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func CreateAccountToken(rtx *redis.Client, token string, account string) (string, error) {
	if err := rtx.Set(token, account, time.Hour*6).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func GetAccountByToken(rtx *redis.Client, token string) (string, error) {
	res, err := rtx.Get(token).Result()
	if err != nil {
		return "", err
	}

	return res, err
}
