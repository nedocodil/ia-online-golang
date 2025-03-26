package users

import (
	"encoding/json"
	"ia-online-golang/internal/domain/services"
	"ia-online-golang/internal/http/responses"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	log         *logrus.Logger
	validator   *validator.Validate
	userService services.UserService
}

// New создаёт новый экземпляр AuthController
func New(log *logrus.Logger, validator *validator.Validate, userService services.UserService) *UserController {
	return &UserController{
		log:         log,
		validator:   validator,
		userService: userService,
	}
}

func (u UserController) Users(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		responses.MethodNotAllowed(w)
		return
	}

	users, err := u.userService.Users(r.Context())
	if err != nil {
		responses.ServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
