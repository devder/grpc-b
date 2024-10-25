package util

import "github.com/go-playground/validator/v10"

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	// check if value is a string
	if cur, ok := fl.Field().Interface().(string); ok {
		return isSupportedCurrency(cur)
	}
	return false
}
