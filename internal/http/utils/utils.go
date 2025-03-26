package utils

import (
	"fmt"
	"ia-online-golang/internal/http/responses"
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

// Кастомный обработчик для несуществующих маршрутов
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	responses.SendError(w, http.StatusNotFound, "По таким путям мы не работаем")
}
