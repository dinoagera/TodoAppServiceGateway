package createtask

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	createtaskrequest "todo-service/internal/models/createTask"
)

type CreateTask interface {
	CreateTask(ctx context.Context, title string, description string) (int64, string, error)
}

func New(log *slog.Logger, createtask CreateTask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Info("user use not allowed method")
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		var request createtaskrequest.CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Info("failed decode")
			http.Error(w, "Error reading JSON", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		id, _, err := createtask.CreateTask(ctx, request.Title, request.Description)
		if err != nil {
			log.Error("failed to created task,", "err:", err)
			http.Error(w, "Error create task", http.StatusBadRequest)
			return
		}
		response := struct {
			ID int64 `json:"id"`
		}{
			ID: id,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
