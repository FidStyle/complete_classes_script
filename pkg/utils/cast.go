package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func CastSliceToString(any interface{}) string {
	res := fmt.Sprintln(any)
	return res[1 : len(res)-2]
}

func CastStringToSlice(any string) []string {
	return strings.Split(any, " ")
}

func FillSameField(source interface{}, dest interface{}) {
	valSource := reflect.ValueOf(source)
	valDest := reflect.ValueOf(dest)

	if valSource.Kind() == reflect.Ptr {
		valSource = valSource.Elem()
	}

	for i := 0; i < valSource.NumField(); i++ {
		fieldSource := valSource.Field(i)
		fieldDest := valDest.Elem().FieldByName(valSource.Type().Field(i).Name)

		if fieldDest.IsValid() && fieldDest.CanSet() && fieldSource.Kind() == fieldDest.Kind() {
			fieldDest.Set(fieldSource)
		}
	}
}
