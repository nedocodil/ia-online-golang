package validator

import (
	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/http/validator/validations"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	v := validator.New()

	// Кастомные валидации
	v.RegisterValidation("complexpassword", validations.PasswordValidation)
	v.RegisterValidation("atLeastOneService", validations.AtLeastOneServiceEnabled)
	v.RegisterStructValidation(validations.NewPasswordStructValidation, dto.NewPasswordDTO{})

	return v
}
