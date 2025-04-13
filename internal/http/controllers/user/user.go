package users

import (
	"encoding/json"
	"errors"
	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/http/context_keys"
	"ia-online-golang/internal/http/responses"
	"ia-online-golang/internal/services/user"
	"ia-online-golang/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	log         *logrus.Logger
	validator   *validator.Validate
	UserService user.UserServiceI
}

type UserControllerI interface {
	User(w http.ResponseWriter, r *http.Request)
	Users(w http.ResponseWriter, r *http.Request)
	EditUser(w http.ResponseWriter, r *http.Request)
}

// New создаёт новый экземпляр AuthController
func New(log *logrus.Logger, validator *validator.Validate, userService user.UserServiceI) *UserController {
	return &UserController{
		log:         log,
		validator:   validator,
		UserService: userService,
	}
}

func (u UserController) User(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		responses.MethodNotAllowed(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	user_id := r.URL.Path[len("/api/v1/user/"):]
	num, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		responses.InvalidRequest(w)
		return
	}

	user, err := u.UserService.UserById(r.Context(), num)
	if err != nil {
		responses.ServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u UserController) Users(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		responses.MethodNotAllowed(w)
		return
	}

	users, err := u.UserService.Users(r.Context())
	if err != nil {
		responses.ServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (u UserController) EditUser(w http.ResponseWriter, r *http.Request) {
	const op = "UserController.EditUser"

	u.log.Debugf("%s: start", op)

	if r.Method != http.MethodPut {
		u.log.Infof("%s: method not allowed. method: %s", op, r.Method)

		w.Header().Set("Allow", http.MethodPut)
		responses.MethodNotAllowed(w)
		return
	}

	u.log.Debugf("%s: method allowed", op)

	var userDTO dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		u.log.Infof("%s: decode error", op)

		responses.InvalidRequest(w)
		return
	}

	u.log.Debugf("%s: decode completed", op)

	err := u.validator.Struct(userDTO)
	if err != nil {
		u.log.Infof("%s: validation error", op)

		responses.ValidationError(w, utils.FormatValidationErrors(err))
		return
	}

	u.log.Debugf("%s: validation completed", op)

	userRolesValue := r.Context().Value(context_keys.UserRoleKey)
	userRoles, ok := userRolesValue.([]string)
	if !ok {
		u.log.Errorf("%s: user role not received", op)

		responses.ServerError(w)
		return
	}

	u.log.Debugf("%s: roles are received", op)

	if userDTO.ID != nil && !utils.Contains(userRoles, "manager") {
		u.log.Infof("%s: forbidden", op)

		responses.Forbidden(w)
		return
	}

	u.log.Debugf("%s: rights checked", op)

	err = u.UserService.EditUser(r.Context(), userDTO)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			u.log.Infof("%s: %v", op, err)

			responses.UserNotFound(w)
			return
		}

		u.log.Errorf("%s: %v", op, err)

		responses.ServerError(w)
		return
	}

	u.log.Debugf("%s: user updated", op)

	w.WriteHeader(http.StatusNoContent)
}
