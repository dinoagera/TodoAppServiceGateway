package donetask

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	donetaskrequest "todo-service/internal/models/doneTask"
)

type DoneTask interface {
	DoneTask(ctx context.Context, id int64) (string, error)
}

func New(log *slog.Logger, donetask DoneTask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Info("user use not allowed method")
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		var request donetaskrequest.DoneTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Info("failed decode")
			http.Error(w, "Error reading JSON", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		message, err := donetask.DoneTask(ctx, request.Id)
		if err != nil {
			log.Info("failed to done task")
			http.Error(w, "Error done task", http.StatusBadRequest)
			return
		}
		response := struct {
			Message string `json:"message"`
		}{
			Message: message,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
