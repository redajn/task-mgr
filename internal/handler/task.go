package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/redajn/task-mgr/internal/domain"
)

type TaskHandler struct {
	svc domain.TaskService
}

func NewTaskHandler(svc domain.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	info, _ := TokenInfoFromContext(r.Context())

	filter := domain.TaskFilter{UserID: info.UserID}

	if doneStr := r.URL.Query().Get("done"); doneStr != "" {
		done, err := strconv.ParseBool(doneStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid 'done' param"})
			return
		}
		filter.Done = &done
	}

	tasks, err := h.svc.ListTask(r.Context(), filter)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	info, _ := TokenInfoFromContext(r.Context())

	var input domain.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	input.UserID = info.UserID

	task, err := h.svc.CreateTask(r.Context(), input)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	info, _ := TokenInfoFromContext(r.Context())

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid task id"})
		return
	}

	task, err := h.svc.GetTask(r.Context(), info.UserID, id)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	info, _ := TokenInfoFromContext(r.Context())

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid task id"})
		return
	}

	var input domain.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	task, err := h.svc.UpdateTask(r.Context(), info.UserID, id, input)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	info, _ := TokenInfoFromContext(r.Context())

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid task id"})
		return
	}

	if err := h.svc.DeleteTask(r.Context(), info.UserID, id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
