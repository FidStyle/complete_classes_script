package order

import (
	"reflect"
	"strings"
	"time"
)

type Order struct {
	// unencrypted password
	ID                  int
	Pw                  string
	Account             string
	PublicRandom        int    `gorm:"column:public_random"`
	ProfessionalRandom  int    `gorm:"column:professional_random"`
	SpecifyPublic       string `gorm:"column:specify_public"`
	SpecifyProfessional string `gorm:"column:specify_professional"`
	B_n                 float64
	F_n                 float64
	A0_n                float64
	A                   float64
	B                   float64
	C                   float64
	D                   float64
	E                   float64
	F                   float64
	A0                  float64
	CreatedAt           time.Time `gorm:"column:created_at"`
	SuccessAt           time.Time `gorm:"column:success_at"`
	Condition           string    `gorm:"column:condition"`
	Creater             string    `gorm:"column:creater"`
}

func (obj *Order) GetScoreByKind(kind string) float64 {
	var score float64 = 0
	val := reflect.ValueOf(obj)
	// typ := reflect.TypeOf(order)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		// typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Float64 && strings.ToLower(val.Type().Field(i).Name) == kind {
			score = field.Float()
			break
		}
	}

	return score
}

var XFSlice = []string{"A", "B", "C", "D", "E", "F", "A0", "B_n", "F_n", "A0_n", "SpecifyPublic", "SpecifyProfessional"}

func (obj *Order) FindUnZeroScore() []string {
	res := []string{}
	val := reflect.ValueOf(obj)
	// typ := reflect.TypeOf(order)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		// typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.Float64 {
			for _, xf := range XFSlice {
				if xf == val.Type().Field(i).Name && val.Field(i).Float() != 0 {
					res = append(res, strings.ToLower(xf))
				}
			}
		}
		if field.Kind() == reflect.String {
			for _, xf := range XFSlice {
				if xf == val.Type().Field(i).Name && val.Field(i).String() != "" {
					res = append(res, xf)
				}
			}
		}
	}

	return res
}
