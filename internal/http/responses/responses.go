package responses

import (
	"encoding/json"
	"net/http"
)

// errorResponse структура для JSON-ответа
type errorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// sendError отправляет JSON-ответ с кодом ошибки
func SendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{
		Message: message,
		Code:    statusCode,
	})
}

func Ok(w http.ResponseWriter) {
	SendError(w, http.StatusOK, "ok")
}
func MethodNotAllowed(w http.ResponseWriter) {
	SendError(w, http.StatusMethodNotAllowed, "method not allowed")
}
func InvalidRequest(w http.ResponseWriter) {
	SendError(w, http.StatusBadRequest, "invalid request")
}
func UserAlreadyExists(w http.ResponseWriter) {
	SendError(w, http.StatusConflict, "user already exists")
}
func UserNotFound(w http.ResponseWriter) {
	SendError(w, http.StatusNotFound, "user not found")
}
func ReferralNotFound(w http.ResponseWriter) {
	SendError(w, http.StatusConflict, "referral not found")
}
func ServerError(w http.ResponseWriter) {
	SendError(w, http.StatusInternalServerError, "error server")
}
func ActivationLinkNotExists(w http.ResponseWriter) {
	SendError(w, http.StatusNotFound, "activation link does not exist")
}
func ActivationLinkExpired(w http.ResponseWriter) {
	SendError(w, http.StatusGone, "activation link expired")
}
func WrongPassword(w http.ResponseWriter) {
	SendError(w, http.StatusUnauthorized, "wrong password")
}
func UserNotActivated(w http.ResponseWriter) {
	SendError(w, http.StatusUnauthorized, "user not activated. activation link sent to email")
}
func ValidationError(w http.ResponseWriter, err string) {
	SendError(w, http.StatusBadRequest, err)
}
func RefreshTokenNotFound(w http.ResponseWriter) {
	SendError(w, http.StatusUnauthorized, "refresh token not found")
}
func AccessTokenNotFound(w http.ResponseWriter) {
	SendError(w, http.StatusUnauthorized, "access token not found")
}
func InvalidRefreshToken(w http.ResponseWriter) {
	SendError(w, http.StatusUnauthorized, "refresh token invalid")
}
func InvalidAccessToken(w http.ResponseWriter) {
	SendError(w, http.StatusUnauthorized, "access token invalid")
}
func ExpiredAccessToken(w http.ResponseWriter) {
	SendError(w, http.StatusUnauthorized, "access token expired")
}
func InvalidBearerFormat(w http.ResponseWriter) {
	SendError(w, http.StatusUnauthorized, "invalid bearer format")
}
func Forbidden(w http.ResponseWriter) {
	SendError(w, http.StatusForbidden, "not enough rights")
}
func OldPasswordIncorrect(w http.ResponseWriter) {
	SendError(w, http.StatusConflict, "old password is incorrect")
}
func PasswordCodeHasExpired(w http.ResponseWriter) {
	SendError(w, http.StatusGone, "password code has expired")
}
func PasswordCodeIncorrect(w http.ResponseWriter) {
	SendError(w, http.StatusUnauthorized, "password code incorrect")
}
