package user

import "time"

type User struct {
	Account   string
	Pw        string
	CreatedAt time.Time `gorm:"column:created_at"`
}
