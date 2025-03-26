package reg

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10/translations/ru"
)

// Кастомная валидация для пароля
func PasswordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Проверка на минимум 1 цифру
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// Проверка на минимум 1 заглавную букву
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Проверка на минимум 1 спецсимвол
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",<>\./?\\|]`).MatchString(password)

	// Проверка на минимум 8 символов
	return len(password) >= 8 && hasDigit && hasUpper && hasSpecial
}
