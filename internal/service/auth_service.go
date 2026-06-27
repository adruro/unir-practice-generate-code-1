package service

import (
	"errors"
	"net/http"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/taskflow/internal/model"
	"github.com/taskflow/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("taskflow-secret-key-change-in-production")

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(username, email, password string) (*model.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("todos los campos son obligatorios")
	}
	if len(password) < 6 {
		return nil, errors.New("la contraseña debe tener al menos 6 caracteres")
	}

	existing, _ := s.userRepo.FindByUsername(username)
	if existing != nil {
		return nil, errors.New("el nombre de usuario ya existe")
	}

	existing, _ = s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("el email ya está registrado")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (*model.User, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("credenciales inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("credenciales inválidas")
	}

	token, err := generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (int64, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return 0, errors.New("token inválido")
	}

	// Verify signature
	signingInput := parts[0] + "." + parts[1]
	expectedSig := sign([]byte(signingInput))
	if parts[2] != expectedSig {
		return 0, errors.New("token inválido")
	}

	// Decode payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, errors.New("token inválido")
	}

	var claims struct {
		UserID int64 `json:"user_id"`
		Exp    int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, errors.New("token inválido")
	}

	if time.Now().Unix() > claims.Exp {
		return 0, errors.New("token expirado")
	}

	return claims.UserID, nil
}

func (s *AuthService) GetUserFromRequest(r *http.Request) (int64, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, errors.New("no autenticado")
	}
	return s.ValidateToken(cookie.Value)
}

func generateToken(userID int64) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	claims := fmt.Sprintf(`{"user_id":%d,"exp":%d}`, userID, time.Now().Add(72*time.Hour).Unix())
	payload := base64.RawURLEncoding.EncodeToString([]byte(claims))

	signingInput := header + "." + payload
	signature := sign([]byte(signingInput))

	return signingInput + "." + signature, nil
}

func sign(data []byte) string {
	h := hmac.New(sha256.New, jwtSecret)
	h.Write(data)
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
