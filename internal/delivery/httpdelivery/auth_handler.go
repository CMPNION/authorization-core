package httpdelivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"cmpnion.space/internal/domain"
	"cmpnion.space/internal/usecase/auth"
)

type AuthService interface {
	SignUp(ctx context.Context, login, password string) (string, error)
	SignIn(ctx context.Context, login, password string) (string, error)
	GetMe(ctx context.Context, userID int64) (auth.UserView, error)
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	var req authRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	token, err := h.service.SignUp(r.Context(), req.Login, req.Password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{Token: token})
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	var req authRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	token, err := h.service.SignIn(r.Context(), req.Login, req.Password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{Token: token})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "Method not allowed"})
		return
	}

	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	user, err := h.service.GetMe(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrLoginTaken):
		writeJSON(w, http.StatusConflict, errorResponse{Error: "login already taken"})
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid credentials"})
	case errors.Is(err, domain.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
	case errors.Is(err, domain.ErrValidation):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request"})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
	}
}

func decodeJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("invalid json: multiple objects")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
