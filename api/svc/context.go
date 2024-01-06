package svc

import (
	"compete_classes_script/api/config"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Cfg *config.Config
	Tx  *gorm.DB
	Rtx *redis.Client
}

func NewServiceContext(c *config.Config) *ServiceContext {
	db, err := gorm.Open(postgres.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rtx := redis.NewClient(&redis.Options{
		Addr:     c.CacheSource,
		Password: "",
		DB:       0,
	})

	return &ServiceContext{
		Cfg: c,
		Tx:  db,
		Rtx: rtx,
	}
}
