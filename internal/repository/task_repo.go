package repository

import (
	"database/sql"

	"github.com/taskflow/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

type TaskFilter struct {
	UserID   int64
	Category string
	Priority string
	Status   string // "all", "pending", "completed"
}

func (r *TaskRepository) Create(task *model.Task) error {
	result, err := r.db.Exec(
		`INSERT INTO tasks (user_id, title, description, category, priority, completed)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		task.UserID, task.Title, task.Description, task.Category, task.Priority, task.Completed,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = id
	return nil
}

func (r *TaskRepository) Update(task *model.Task) error {
	_, err := r.db.Exec(
		`UPDATE tasks SET title=?, description=?, category=?, priority=?, completed=?, updated_at=CURRENT_TIMESTAMP
		 WHERE id=? AND user_id=?`,
		task.Title, task.Description, task.Category, task.Priority, task.Completed, task.ID, task.UserID,
	)
	return err
}

func (r *TaskRepository) Delete(id, userID int64) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id=? AND user_id=?", id, userID)
	return err
}

func (r *TaskRepository) FindByID(id, userID int64) (*model.Task, error) {
	task := &model.Task{}
	err := r.db.QueryRow(
		`SELECT id, user_id, title, description, category, priority, completed, created_at, updated_at
		 FROM tasks WHERE id=? AND user_id=?`, id, userID,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Category,
		&task.Priority, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepository) FindAll(filter TaskFilter) ([]model.Task, error) {
	query := `SELECT id, user_id, title, description, category, priority, completed, created_at, updated_at
			  FROM tasks WHERE user_id = ?`
	args := []interface{}{filter.UserID}

	if filter.Category != "" && filter.Category != "todas" {
		query += " AND category = ?"
		args = append(args, filter.Category)
	}
	if filter.Priority != "" && filter.Priority != "todas" {
		query += " AND priority = ?"
		args = append(args, filter.Priority)
	}
	if filter.Status == "pending" {
		query += " AND completed = 0"
	} else if filter.Status == "completed" {
		query += " AND completed = 1"
	}

	query += " ORDER BY completed ASC, CASE priority WHEN 'alta' THEN 1 WHEN 'media' THEN 2 WHEN 'baja' THEN 3 END, created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Category,
			&t.Priority, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) ToggleComplete(id, userID int64) error {
	_, err := r.db.Exec(
		`UPDATE tasks SET completed = CASE WHEN completed = 0 THEN 1 ELSE 0 END, updated_at=CURRENT_TIMESTAMP
		 WHERE id=? AND user_id=?`, id, userID,
	)
	return err
}

func (r *TaskRepository) CountByUser(userID int64) (total int, completed int, err error) {
	err = r.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE user_id=?", userID).Scan(&total)
	if err != nil {
		return
	}
	err = r.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE user_id=? AND completed=1", userID).Scan(&completed)
	return
}
