package model

import (
	"errors"
	"fmt"
	validate "github.com/go-playground/validator/v10"
)

var (
	Validator = validate.New()
)

func ValidationError(err error) []string {
	var validationErrors validate.ValidationErrors
	errorSlice := make([]string, 0)
	if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
		errorSlice = append(errorSlice, fmt.Sprint(validationErrors[0].Field()+" failed on the "+validationErrors[0].Tag()+" tag"))
	}

	return errorSlice
}
