// internal/http/controllers/auth/auth.go
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/http/context_keys"
	"ia-online-golang/internal/http/responses"
	"ia-online-golang/internal/utils"

	"ia-online-golang/internal/services/auth"
	"ia-online-golang/internal/services/passwordcode"
	"ia-online-golang/internal/services/token"
	"ia-online-golang/internal/services/user"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	log         *logrus.Logger
	validator   *validator.Validate
	AuthService auth.AuthServiceI
}

type AuthControllerI interface {
	Registration(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Activation(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	NewPassword(w http.ResponseWriter, r *http.Request)
	// SendPasswordCode(w http.ResponseWriter, r *http.Request)
	// RecoverPassword(w http.ResponseWriter, r *http.Request)
}

// New создаёт новый экземпляр AuthController
func New(log *logrus.Logger, validator *validator.Validate, authService auth.AuthServiceI) *AuthController {
	return &AuthController{
		AuthService: authService,
		log:         log,
		validator:   validator,
	}
}

func (a *AuthController) Registration(w http.ResponseWriter, r *http.Request) {
	const op = "Controller.Registration"

	a.log.Debugf("%s: start", op)

	if r.Method != http.MethodPost {
		a.log.Infof("%s: method not allowed", op)

		w.Header().Set("Allow", http.MethodPost)
		responses.MethodNotAllowed(w)
		return
	}

	a.log.Debugf("%s: method is correct", op)

	var dto dto.RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		a.log.Infof("%s: %v", op, err)

		responses.InvalidRequest(w)
		return
	}

	// Валидируем данные
	if err := a.validator.Struct(dto); err != nil {
		a.log.Infof("%s: %v", op, err)

		responses.ValidationError(w, utils.FormatValidationErrors(err))
		return
	}

	a.log.Debugf("%s: validation completed", op)

	// Регистрируем пользователя
	tokens, err := a.AuthService.RegistrationUser(r.Context(), dto)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			a.log.Infof("%s: %v", op, err)

			responses.UserAlreadyExists(w)
			return
		}
		if errors.Is(err, auth.ErrReferralIdNotFound) {
			a.log.Infof("%s: %v", op, err)

			responses.ReferralNotFound(w)
			return
		}

		a.log.Errorf("%s: %v", op, err)

		responses.ServerError(w)
		return
	}

	a.log.Debugf("%s: user registration", op)

	// Создаем cookie с токеном
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   3600 * 24 * 30, // 30 дней
	}

	http.SetCookie(w, cookie)

	a.log.Debugf("%s: tokens send", op)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
func (a *AuthController) Activation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		responses.MethodNotAllowed(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	activation_id := r.URL.Path[len("/api/v1/auth/activation/"):]

	err := a.AuthService.ActivationUser(r.Context(), activation_id)
	if err != nil {
		if errors.Is(err, auth.ErrActiveLinkNotExists) {
			responses.ActivationLinkNotExists(w)
			return
		}

		if errors.Is(err, auth.ErrActiveLinkExpired) {
			responses.ActivationLinkExpired(w)
			return
		}

		responses.ServerError(w)
		return
	}

	http.Redirect(w, r, "/auth/test", http.StatusSeeOther)
}
func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	op := "Controller.Login"

	a.log.Debugf("%s: start", op)

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		a.log.Infof("%s: method not allowed", op)
		responses.MethodNotAllowed(w)
		return
	}

	a.log.Debugf("%s: method is correct", op)

	var dto dto.LoginUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		a.log.Infof("%s: %v", op, err)

		responses.InvalidRequest(w)
		return
	}

	// Валидируем данные
	if err := a.validator.Struct(dto); err != nil {
		a.log.Infof("%s: %v", op, err)

		responses.ValidationError(w, utils.FormatValidationErrors(err))
		return
	}

	a.log.Debugf("%s: validation completed", op)

	tokens, err := a.AuthService.LoginUser(r.Context(), dto)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			a.log.Infof("%s: user not found", op)
			responses.UserNotFound(w)
			return
		}

		if errors.Is(err, auth.ErrIncorrectPassword) {
			a.log.Infof("%s: password incorrect", op)
			responses.WrongPassword(w)
			return
		}

		if errors.Is(err, user.ErrUserNotActivated) {
			a.log.Infof("%s: user not activated", op)
			responses.UserNotActivated(w)
			return
		}

		a.log.Errorf("%s: server error: %v", op, err)
		responses.ServerError(w)
		return
	}

	a.log.Debugf("%s: token created", op)

	// Создаем cookie с токеном
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   3600 * 24 * 30, // 30 дней
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	a.log.Debugf("%s: refresh token add from cookie", op)

	w.Header().Set("Content-Type", "application/json")
	a.log.Infof("%s: tokens send", op)
	json.NewEncoder(w).Encode(tokens)
}
func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		responses.MethodNotAllowed(w)
		return
	}

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		responses.RefreshTokenNotFound(w)
		return
	}

	err = a.AuthService.LogoutUser(r.Context(), refreshToken.Value)
	if err != nil {
		responses.ServerError(w)
		return
	}

	// Создаем cookie с токеном
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   -1, // 30 дней
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}
func (a *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	op := "AuthController.Refresh"

	a.log.Debugf("%s: %v", op, r.RemoteAddr)

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		responses.MethodNotAllowed(w)
		return
	}

	a.log.Debugf("%s: method is correct", op)

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		a.log.Infof("%s: refresh token not found", op)

		responses.RefreshTokenNotFound(w)
		return
	}

	a.log.Debugf("%s: token received %v", op, refreshToken.Value)

	tokens, err := a.AuthService.RefreshUserTokens(r.Context(), refreshToken.Value)
	if err != nil {
		if errors.Is(err, token.ErrInvalidRefreshToken) {
			a.log.Infof("%s: invalid refresh token", op)
			responses.InvalidRefreshToken(w)
			return
		}

		if errors.Is(err, token.ErrRefreshTokenNotExists) {
			a.log.Infof("%s: user not activated", op)
			responses.RefreshTokenNotFound(w)
			return
		}

		a.log.Errorf("%s: %v", op, err)
		responses.ServerError(w)
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   3600 * 24 * 30, // 30 дней
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	a.log.Debugf("%s: cookie updated", op)

	w.Header().Set("Content-Type", "application/json")
	a.log.Infof("%s: tokens send", op)
	json.NewEncoder(w).Encode(tokens)
}
func (a *AuthController) NewPassword(w http.ResponseWriter, r *http.Request) {
	op := "AuthController.NewPassword"

	a.log.Debugf("%s: start", op)

	if r.Method != http.MethodPost {
		a.log.Infof("%s: method not allowed. method: %s", op, r.Method)

		w.Header().Set("Allow", http.MethodPost)
		responses.MethodNotAllowed(w)
		return
	}

	a.log.Debugf("%s: method is correct", op)

	var dto dto.NewPasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		a.log.Infof("%s: invalid request", op)

		responses.InvalidRequest(w)
		return
	}

	// Валидируем данные
	if err := a.validator.Struct(dto); err != nil {
		a.log.Infof("%s: invalid request", op)

		responses.ValidationError(w, utils.FormatValidationErrors(err))
		return
	}

	a.log.Debugf("%s: validation completed", op)

	userIDValue := r.Context().Value(context_keys.UserIDKey)
	userID, ok := userIDValue.(int64)
	if !ok {
		a.log.Errorf("%s: id user not received", op)

		responses.ServerError(w)
		return
	}

	a.log.Debugf("%s: id received", op)

	err := a.AuthService.ChangingPassword(r.Context(), dto, userID)
	if err != nil {
		if errors.Is(err, auth.ErrIncorrectOldPassword) {
			a.log.Infof("%s: old password incorrect", op)

			responses.OldPasswordIncorrect(w)
			return
		}
		a.log.Errorf("%s: %v", op, err)

		responses.ServerError(w)
		return
	}

	a.log.Debugf("%s: password changing", op)

	responses.Ok(w)
}
func (a *AuthController) SendNewPassword(w http.ResponseWriter, r *http.Request) {
	const op = "AuthController.SendNewPassword"

	a.log.Debugf("%s: %v", op, r.RemoteAddr)

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		responses.MethodNotAllowed(w)
		return
	}

	a.log.Debugf("%s: method is correct", op)

	var dto dto.SendNewPasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		a.log.Infof("%s: %v", op, err)

		responses.InvalidRequest(w)
		return
	}

	if err := a.validator.Struct(dto); err != nil {
		a.log.Infof("%s: %v", op, err)

		responses.ValidationError(w, utils.FormatValidationErrors(err))
		return
	}

	a.log.Debugf("%s: validation completed", op)

	err := a.AuthService.RecoverPassword(r.Context(), dto.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			responses.UserNotFound(w)
			return
		}

		responses.ServerError(w)
		return
	}

	responses.Ok(w)
}
func (a *AuthController) SendPasswordCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		responses.MethodNotAllowed(w)
		return
	}

	var dto dto.SendPasswordCodeDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		responses.InvalidRequest(w)
		return
	}

	if err := a.validator.Struct(dto); err != nil {
		responses.ValidationError(w, utils.FormatValidationErrors(err))
		return
	}

	err := a.AuthService.RecoverPassword(r.Context(), dto.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			responses.UserNotFound(w)
			return
		}

		responses.ServerError(w)
		return
	}

	responses.Ok(w)
}
func (a *AuthController) RecoverPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		responses.MethodNotAllowed(w)
		return
	}

	var dto dto.RecoverPasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		responses.InvalidRequest(w)
		return
	}

	if err := a.validator.Struct(dto); err != nil {
		responses.ValidationError(w, utils.FormatValidationErrors(err))
		return
	}

	err := a.AuthService.NewPassword(r.Context(), dto)
	if err != nil {
		if errors.Is(err, passwordcode.ErrPasswordCodeIsNotFound) || errors.Is(err, passwordcode.ErrPasswordCodeIncorrect) {
			responses.PasswordCodeIncorrect(w)
			return
		}

		if errors.Is(err, passwordcode.ErrPasswordCodeHasExpired) {
			responses.PasswordCodeHasExpired(w)
			return
		}

		fmt.Println(err)

		responses.ServerError(w)
		return
	}

	responses.Ok(w)
}
