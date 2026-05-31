package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FiberStructValidator struct {
	validate *validator.Validate
}

func NewFiberStructValidator() *FiberStructValidator {
	return &FiberStructValidator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *FiberStructValidator) Validate(out any) error {
	if err := v.validate.Struct(out); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMessages []string
			for _, fieldErr := range validationErrors {
				errMessages = append(errMessages,
					fmt.Sprintf("field '%s' failed on '%s'", fieldErr.Field(), fieldErr.Tag()))
			}
			return fmt.Errorf("validation errors: %s", strings.Join(errMessages, "; "))
		}
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
