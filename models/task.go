package models

import (
    "time"
)

type Task struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
}

type TaskRequest struct {
    Title     string `json:"title" binding:"required"`
    Completed bool   `json:"completed"`
}

type TaskStats struct {
    Total     int `json:"total"`
    Completed int `json:"completed"`
    Active    int `json:"active"`
}

type PaginatedResponse struct {
    Tasks  []Task `json:"tasks"`
    Total  int    `json:"total"`
    Limit  int    `json:"limit"`
    Offset int    `json:"offset"`
}