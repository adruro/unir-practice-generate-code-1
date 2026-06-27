package service

import (
	"errors"

	"github.com/taskflow/internal/model"
	"github.com/taskflow/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) Create(userID int64, title, description, category, priority string) (*model.Task, error) {
	if title == "" {
		return nil, errors.New("el título es obligatorio")
	}

	if !isValidCategory(category) {
		category = "personal"
	}
	if !isValidPriority(priority) {
		priority = "media"
	}

	task := &model.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Category:    model.Category(category),
		Priority:    model.Priority(priority),
		Completed:   false,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Update(id, userID int64, title, description, category, priority string) error {
	task, err := s.taskRepo.FindByID(id, userID)
	if err != nil {
		return errors.New("tarea no encontrada")
	}

	if title != "" {
		task.Title = title
	}
	task.Description = description
	if isValidCategory(category) {
		task.Category = model.Category(category)
	}
	if isValidPriority(priority) {
		task.Priority = model.Priority(priority)
	}

	return s.taskRepo.Update(task)
}

func (s *TaskService) Delete(id, userID int64) error {
	return s.taskRepo.Delete(id, userID)
}

func (s *TaskService) ToggleComplete(id, userID int64) error {
	return s.taskRepo.ToggleComplete(id, userID)
}

func (s *TaskService) GetAll(userID int64, category, priority, status string) ([]model.Task, error) {
	filter := repository.TaskFilter{
		UserID:   userID,
		Category: category,
		Priority: priority,
		Status:   status,
	}
	return s.taskRepo.FindAll(filter)
}

func (s *TaskService) GetByID(id, userID int64) (*model.Task, error) {
	return s.taskRepo.FindByID(id, userID)
}

func (s *TaskService) GetStats(userID int64) (total int, completed int, err error) {
	return s.taskRepo.CountByUser(userID)
}

func isValidCategory(c string) bool {
	return c == "trabajo" || c == "personal" || c == "estudio"
}

func isValidPriority(p string) bool {
	return p == "alta" || p == "media" || p == "baja"
}
