package infra

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return Validator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *Validator) Validate(body interface{}) []string {
	err := h.validate.Struct(body)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := []string{}

		for _, validationError := range validationErrors {
			tag := validationError.Tag()
			param := validationError.Param()
			field := strings.ToLower(validationError.Field()[:1]) + validationError.Field()[1:]

			switch tag {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", field))
			case "uuid4":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be uuidv4", field))
			case "gte":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be greater than or equal to %s", field, param))
			}
		}

		return errorMessages
	}

	return []string{}
}
