package model

import "time"

type Priority string

const (
	PriorityHigh   Priority = "alta"
	PriorityMedium Priority = "media"
	PriorityLow    Priority = "baja"
)

type Category string

const (
	CategoryWork     Category = "trabajo"
	CategoryPersonal Category = "personal"
	CategoryStudy    Category = "estudio"
)

type Task struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    Category  `json:"category"`
	Priority    Priority  `json:"priority"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
