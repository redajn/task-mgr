package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/redajn/task-mgr/internal/domain"
)

type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrTaskNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{
			Error: "task not found",
			Code:  "TASK_NOT_FOUND",
		})
	case errors.Is(err, domain.ErrTaskForbidden):
		writeJSON(w, http.StatusForbidden, errorResponse{
			Error: "title cannot be empty",
			Code:  "FORBIDDEN",
		})
	case errors.Is(err, domain.ErrTaskTitleEmpty):
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "title cannot be empty",
			Code:  "VALIDATION_ERROR",
		})
	default:
		slog.Error("internal error", "error", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "internal sercer error",
		})
	}
}
