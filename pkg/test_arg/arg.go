package testarg

import (
	"compete_classes_script/api/config"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Tx() *gorm.DB {
	c := config.NewConfig("/home/sagayosa/compete_classes_script/config.yaml")

	db, err := gorm.Open(postgres.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func Rtx() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}
