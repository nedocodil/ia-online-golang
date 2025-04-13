package utils

import (
	"crypto/rand"
	"fmt"
	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/http/responses"
	"ia-online-golang/internal/models"
	"math/big"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) string {
	var messages []string
	for _, e := range err.(validator.ValidationErrors) {
		messages = append(messages, fmt.Sprintf("Field '%s' is invalid", e.Field()))
	}
	return strings.Join(messages, ", ")
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	responses.SendError(w, http.StatusNotFound, "По таким путям мы не работаем")
}

func UserToDTO(user models.User) dto.UserDTO {
	return dto.UserDTO{
		ID:           &user.ID,
		Roles:        user.Roles,
		ReferralCode: user.ReferralCode,
		Email:        user.Email,
		Name:         user.Name,
		PhoneNumber:  user.PhoneNumber,
		City:         user.City,
		Telegram:     user.Telegram,
	}
}

func DtoToUser(user dto.UserDTO) models.User {
	return models.User{
		ID:           derefInt64(user.ID),
		Roles:        user.Roles,
		ReferralCode: user.ReferralCode,
		Email:        user.Email,
		Name:         user.Name,
		PhoneNumber:  user.PhoneNumber,
		City:         user.City,
		Telegram:     user.Telegram,
	}
}

func UserToReferralDTO(user models.User, status bool) dto.ReferralDTO {
	return dto.ReferralDTO{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		City:        user.City,
		Active:      status,
	}
}

func GeneratePasswordCode(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func derefInt64(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}
