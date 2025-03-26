// internal/http/controllers/auth/auth.go
package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"ia-online-golang/internal/domain/services"
	"ia-online-golang/internal/http/responses"
	"ia-online-golang/internal/http/utils"
	"ia-online-golang/internal/interfaces/dto"

	ServiceAuthErrors "ia-online-golang/internal/interfaces/errors/services"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	log         *logrus.Logger
	AuthService services.AuthService
	validator   *validator.Validate
}

// New создаёт новый экземпляр AuthController
func New(authService services.AuthService, log *logrus.Logger, validator *validator.Validate) *AuthController {
	return &AuthController{
		AuthService: authService,
		log:         log,
		validator:   validator,
	}
}

func (a AuthController) Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		responses.MethodNotAllowed(w)
		return
	}

	var dto dto.RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		responses.InvalidRequest(w)
		return
	}

	// Валидируем данные
	if err := a.validator.Struct(dto); err != nil {
		responses.ValidationError(w, utils.FormatValidationErrors(err))
		return
	}

	// Регистрируем пользователя
	tokens, _, err := a.AuthService.RegistrationUser(r.Context(), dto)
	if err != nil {
		if errors.Is(err, ServiceAuthErrors.ErrUserAlreadyExists) {
			responses.UserAlreadyExists(w)
			return
		}
		if errors.Is(err, ServiceAuthErrors.ErrReferralAlreadyUsed) {
			responses.ReferralNotUsed(w)
			return
		}
		responses.ServerError(w)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
func (a AuthController) Activation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		responses.MethodNotAllowed(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	activation_id := r.URL.Path[len("/api/v1/auth/activation/"):]

	err := a.AuthService.ActivationUser(r.Context(), activation_id)
	if err != nil {
		if errors.Is(err, ServiceAuthErrors.ErrActiveLinkNotExists) {
			responses.ActivationLinkNotExists(w)
			return
		}

		if errors.Is(err, ServiceAuthErrors.ErrActiveLinkExpired) {
			responses.ActivationLinkExpired(w)
			return
		}

		responses.ServerError(w)
		return
	}

	http.Redirect(w, r, "/auth/test", http.StatusSeeOther)
}
func (a AuthController) Login(w http.ResponseWriter, r *http.Request) {
	op := "Controller.Login"

	a.log.Debugf("%s: start", op)

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		a.log.Info("method not allowed")
		responses.MethodNotAllowed(w)
		return
	}

	a.log.Debugf("%s: method is correct", op)

	var dto dto.LoginUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		a.log.Infof("%s: Validation failed", op)
		responses.InvalidRequest(w)
		return
	}

	a.log.Debugf("%s: validation completed", op)

	tokens, err := a.AuthService.LoginUser(r.Context(), dto)
	if err != nil {
		if errors.Is(err, ServiceAuthErrors.ErrUserNotFound) {
			a.log.Infof("%s: user not found", op)
			responses.UserNotFound(w)
			return
		}

		if errors.Is(err, ServiceAuthErrors.ErrIncorrectPassword) {
			a.log.Infof("%s: password incorrect", op)
			responses.WrongPassword(w)
			return
		}

		if errors.Is(err, ServiceAuthErrors.ErrUserNotActivated) {
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
func (a AuthController) Logout(w http.ResponseWriter, r *http.Request) {
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
func (a AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	op := "AuthController.Refresh"

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		responses.MethodNotAllowed(w)
		return
	}

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		responses.RefreshTokenNotFound(w)
		return
	}

	tokens, err := a.AuthService.RefreshUserToken(r.Context(), refreshToken.Value)
	if err != nil {
		if errors.Is(err, ServiceAuthErrors.ErrInvalidRefreshToken) {
			a.log.Infof("%s: invalid refresh token", op)
			responses.InvalidRefreshToken(w)
			return
		}

		if errors.Is(err, ServiceAuthErrors.ErrRefreshTokenNotExists) {
			a.log.Infof("%s: user not activated", op)
			responses.RefreshTokenNotFound(w)
			return
		}

		responses.ServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	a.log.Infof("%s: tokens send", op)
	json.NewEncoder(w).Encode(tokens)
}
