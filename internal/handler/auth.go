package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/redajn/task-mgr/internal/domain"
)

type AuthHandler struct {
	svc domain.AuthService
}

func NewAuthHandler(svc domain.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input domain.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	user, err := h.svc.Register(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			writeJSON(w, http.StatusConflict, errorResponse{
				Error: "user already exists",
				Code:  "USER_ALREDY_EXISTS",
			})
		default:
			handleError(w, err)
		}
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input domain.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	token, err := h.svc.Login(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			writeJSON(w, http.StatusUnauthorized, errorResponse{
				Error: "invalid email or password",
				Code:  "INVALID_CREDENTIALS",
			})
		default:
			handleError(w, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing token"})
		return
	}

	if err := h.svc.Logout(r.Context(), token); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func extractToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}
