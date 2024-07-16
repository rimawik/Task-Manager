package models

// Task represents a task in the system
// @Description Task represents a task in the system
type Task struct {
	ID          int     `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	CreateTime  *int64  `json:"create_time"`
	UpdateTime  *int64  `json:"update_time"`
	Deadline    *int64  `json:"deadline"`
}
