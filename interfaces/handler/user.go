package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"github.com/tusmasoma/go-tech-dojo/usecase"
)

type UserHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	uuc usecase.UserUseCase
}

func NewUserHandler(uuc usecase.UserUseCase) UserHandler {
	return &userHandler{
		uuc: uuc,
	}
}

type GetUserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Coins     int    `json:"coins"`
	HighScore int    `json:"high_score"`
}

func (uh *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := uh.uuc.GetUser(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Coins:     user.Coins,
		HighScore: user.HighScore,
	}); err != nil {
		http.Error(w, "Failed to encode user to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody CreateUserRequest
	defer r.Body.Close()
	if !uh.isValidCreateUserRequest(r.Body, &requestBody) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := uh.uuc.CreateUserAndToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (uh *userHandler) isValidCreateUserRequest(body io.ReadCloser, requestBody *CreateUserRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body: %v", err)
		return false
	}
	if requestBody.Email == "" || requestBody.Password == "" {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}

type UpdateUserRequest struct {
	Coins     int `json:"coins"`
	HighScore int `json:"high_score"`
}

type UpdateUserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Coins     int    `json:"coins"`
	HighScore int    `json:"high_score"`
}

func (uh *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody UpdateUserRequest
	defer r.Body.Close()
	if !uh.isValidUpdateUserRequest(r.Body, &requestBody) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := uh.uuc.UpdateUser(ctx, requestBody.Coins, requestBody.HighScore)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(UpdateUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Coins:     user.Coins,
		HighScore: user.HighScore,
	}); err != nil {
		http.Error(w, "Failed to encode user to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (uh *userHandler) isValidUpdateUserRequest(body io.ReadCloser, requestBody *UpdateUserRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body: %v", err)
		return false
	}
	if requestBody.Coins < 0 || requestBody.HighScore < 0 {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}
