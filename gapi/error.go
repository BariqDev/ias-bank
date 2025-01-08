package gapi

import (
	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

var ErrorMessages = map[string]string{
	"username":  "must be alphanumeric and required",
	"email":     "must be a valid email address",
	"password":  "must be at least 6 characters long",
	"full_name": "must be required",
}

func fieldsViolations(errs []validator.FieldError) []*errdetails.BadRequest_FieldViolation {
	if len(errs) == 0 {
		return nil
	}

	var violations []*errdetails.BadRequest_FieldViolation
	var errorMsg string
	
	for _, err := range errs {

		if message, ok := ErrorMessages[err.Field()]; ok {
			errorMsg = message
		} else {
			errorMsg = err.Error()
		}
		var errField = &errdetails.BadRequest_FieldViolation{
			Field:       err.Field(),
			Description: errorMsg,
		}
		violations = append(violations, errField)
	}

	return violations
}
