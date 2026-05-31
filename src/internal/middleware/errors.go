package middleware

import "fmt"

type ParamsValidationError struct {
	err error
}

func (validationError ParamsValidationError) Error() string {
	return fmt.Sprintf("failed to bind parameters: %s", validationError.err.Error())
}

func (validationError ParamsValidationError) StatusCode() int {
	return 400
}

type BodyValidationError struct {
	err error
}

func (validationError BodyValidationError) Error() string {
	return fmt.Sprintf("invalid request body: %s", validationError.err.Error())
}

func (validationError BodyValidationError) StatusCode() int {
	return 400
}
