package validations

import (
	"crypto/rand"
	"fmt"
	"ia-online-golang/internal/dto"
	"math/big"
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

func AtLeastOneServiceEnabled(fl validator.FieldLevel) bool {
	obj := fl.Parent().Interface().(dto.LeadDTO)
	return obj.IsInternet || obj.IsShipping || obj.IsCleaning
}

func NewPasswordStructValidation(sl validator.StructLevel) {
	dto := sl.Current().Interface().(dto.NewPasswordDTO)

	if dto.OldPassword == dto.NewPassword {
		// Регистрируем ошибку на поле NewPassword, можно изменить сообщение в переводах
		sl.ReportError(dto.NewPassword, "NewPassword", "new_password", "nefield", "OldPassword")
	}
}

func GenerateValidPassword(length int) (string, error) {
	if length < 8 {
		return "", fmt.Errorf("password length must be at least 8 characters")
	}

	digits := "0123456789"
	uppers := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specials := "!@#$%^&*()_+-=[]{};:'\",.<>/?\\|"
	all := digits + uppers + specials + "abcdefghijklmnopqrstuvwxyz"

	// Гарантируем хотя бы по одному символу каждого типа
	password := make([]byte, length)
	var err error

	password[0], err = randomChar(digits)
	if err != nil {
		return "", err
	}
	password[1], err = randomChar(uppers)
	if err != nil {
		return "", err
	}
	password[2], err = randomChar(specials)
	if err != nil {
		return "", err
	}

	// Остальные символы
	for i := 3; i < length; i++ {
		password[i], err = randomChar(all)
		if err != nil {
			return "", err
		}
	}

	// Перемешиваем пароль
	shuffle(password)

	return string(password), nil
}

func randomChar(charset string) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		return 0, err
	}
	return charset[n.Int64()], nil
}

func shuffle(data []byte) {
	for i := range data {
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(data))))
		j := int(jBig.Int64())
		data[i], data[j] = data[j], data[i]
	}
}
