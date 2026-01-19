package services

import (
    "strings"
    
    "GOproject/models"
    "GOproject/repository"
)

type TaskService struct {
    repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

func (s *TaskService) GetAllTasks(status, search string, limit, offset int) (*models.PaginatedResponse, error) {
    tasks, total, err := s.repo.FindAll(status, search, limit, offset)
    if err != nil {
        return nil, err
    }
    
    return &models.PaginatedResponse{
        Tasks:  tasks,
        Total:  total,
        Limit:  limit,
        Offset: offset,
    }, nil
}

func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
    return s.repo.FindByID(id)
}

func (s *TaskService) CreateTask(taskReq *models.TaskRequest) (*models.Task, error) {
    // Валидация
    if err := s.validateTaskRequest(taskReq); err != nil {
        return nil, err
    }
    
    // Создаем задачу
    id, err := s.repo.Create(taskReq)
    if err != nil {
        return nil, err
    }
    
    // Получаем созданную задачу
    task, err := s.repo.FindByID(int(id))
    if err != nil {
        return nil, err
    }
    
    return task, nil
}

func (s *TaskService) UpdateTask(id int, taskReq *models.TaskRequest) (*models.Task, error) {
    // Валидация
    if err := s.validateTaskRequest(taskReq); err != nil {
        return nil, err
    }
    
    // Проверяем существование задачи
    exists, err := s.repo.Exists(id)
    if err != nil {
        return nil, err
    }
    
    if !exists {
        return nil, &NotFoundError{ID: id}
    }
    
    // Обновляем задачу
    if err := s.repo.Update(id, taskReq); err != nil {
        return nil, err
    }
    
    // Получаем обновленную задачу
    return s.repo.FindByID(id)
}

func (s *TaskService) DeleteTask(id int) error {
    // Проверяем существование задачи
    exists, err := s.repo.Exists(id)
    if err != nil {
        return err
    }
    
    if !exists {
        return &NotFoundError{ID: id}
    }
    
    return s.repo.Delete(id)
}

func (s *TaskService) GetStats() (*models.TaskStats, error) {
    return s.repo.GetStats()
}

func (s *TaskService) HealthCheck() error {
    return s.repo.CheckHealth()
}

func (s *TaskService) validateTaskRequest(taskReq *models.TaskRequest) error {
    taskReq.Title = strings.TrimSpace(taskReq.Title)
    
    if taskReq.Title == "" {
        return &ValidationError{Field: "title", Message: "Заголовок задачи не может быть пустым"}
    }
    
    if len(taskReq.Title) > 500 {
        return &ValidationError{Field: "title", Message: "Заголовок слишком длинный (максимум 500 символов)"}
    }
    
    return nil
}

// Custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return e.Message
}

type NotFoundError struct {
    ID int
}

func (e *NotFoundError) Error() string {
    return "Задача не найдена"
}