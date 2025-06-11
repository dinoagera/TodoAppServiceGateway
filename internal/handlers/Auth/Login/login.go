package login

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	loginandregisteruser "todo-service/internal/models/loginAndRegisterUser"
)

type Login interface {
	Login(ctx context.Context, email string, password string) (string, error)
}

func New(log *slog.Logger, login Login) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Info("user use not allowed method")
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		var request loginandregisteruser.LoginAndRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Info("failed decode")
			http.Error(w, "Error reading JSON", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		token, err := login.Login(ctx, request.Email, request.Password)
		if err != nil {
			log.Error("login failed", "error", err)
			http.Error(w, "Error login", http.StatusBadRequest)
			return
		}
		response := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

	}
}
