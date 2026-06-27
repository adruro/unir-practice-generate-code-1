package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/taskflow/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Check if already logged in
	if _, err := h.authService.GetUserFromRequest(r); err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "login.html", nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	_, token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		jsonError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(72 * time.Hour),
	})

	jsonResponse(w, map[string]string{"message": "Login exitoso", "redirect": "/dashboard"})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	_, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Auto login after register
	_, token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		jsonError(w, "Registro exitoso. Por favor inicia sesión.", http.StatusOK)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(72 * time.Hour),
	})

	jsonResponse(w, map[string]string{"message": "Registro exitoso", "redirect": "/dashboard"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
