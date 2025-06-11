package getalltask

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"todo-service/internal/models/task"
)

type GetAllTask interface {
	GetAllTask(ctx context.Context) ([]*task.Task, error)
}

func New(log *slog.Logger, getalltask GetAllTask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Info("user use not allowed method")
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		tasks, err := getalltask.GetAllTask(ctx)
		if err != nil {
			log.Info("failed to get all task")
			http.Error(w, "Error get all task", http.StatusBadRequest)
			return
		}
		response := struct {
			Tasks []*task.Task `json:"tasks"`
		}{
			Tasks: tasks,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
