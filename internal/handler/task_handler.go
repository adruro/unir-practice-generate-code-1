package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/taskflow/internal/model"
	"github.com/taskflow/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
	authService *service.AuthService
}

func NewTaskHandler(taskService *service.TaskService, authService *service.AuthService) *TaskHandler {
	return &TaskHandler{taskService: taskService, authService: authService}
}

func (h *TaskHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	total, completed, _ := h.taskService.GetStats(userID)
	data := map[string]interface{}{
		"Total":     total,
		"Completed": completed,
		"Pending":   total - completed,
	}
	renderTemplate(w, "dashboard.html", data)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		jsonError(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	category := r.URL.Query().Get("category")
	priority := r.URL.Query().Get("priority")
	status := r.URL.Query().Get("status")

	tasks, err := h.taskService.GetAll(userID, category, priority, status)
	if err != nil {
		jsonError(w, "Error al obtener tareas", http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	jsonResponse(w, tasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		jsonError(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Priority    string `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.Create(userID, req.Title, req.Description, req.Category, req.Priority)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		jsonError(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Priority    string `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if err := h.taskService.Update(id, userID, req.Title, req.Description, req.Category, req.Priority); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, map[string]string{"message": "Tarea actualizada"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		jsonError(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.taskService.Delete(id, userID); err != nil {
		jsonError(w, "Error al eliminar tarea", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": "Tarea eliminada"})
}

func (h *TaskHandler) ToggleTask(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		jsonError(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.taskService.ToggleComplete(id, userID); err != nil {
		jsonError(w, "Error al actualizar tarea", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": "Estado actualizado"})
}
