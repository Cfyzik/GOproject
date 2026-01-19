package repository

import (
    "database/sql"
    "fmt"
    "log"
    "strings"
    "time"
    
    "GOproject/models"
    "GOproject/config"
)

type TaskRepository struct {
    db *sql.DB
}

func NewTaskRepository() *TaskRepository {
    return &TaskRepository{db: config.DB}
}

func (r *TaskRepository) FindAll(status, search string, limit, offset int) ([]models.Task, int, error) {
    query := "SELECT id, title, completed, created_at FROM tasks"
    conditions := []string{}
    args := []interface{}{}
    
    // Добавляем фильтр по статусу
    switch status {
    case "completed":
        conditions = append(conditions, "completed = ?")
        args = append(args, true)
    case "active":
        conditions = append(conditions, "completed = ?")
        args = append(args, false)
    }
    
    // Добавляем поиск по заголовку
    if search != "" {
        conditions = append(conditions, "title LIKE ?")
        args = append(args, "%"+search+"%")
    }
    
    // Объединяем условия WHERE
    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }
    
    // Добавляем сортировку и пагинацию
    query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
    args = append(args, limit, offset)
    
    // Выполняем запрос
    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("ошибка выполнения запроса: %v", err)
    }
    defer rows.Close()
    
    // Собираем результаты
    tasks := []models.Task{}
    for rows.Next() {
        var task models.Task
        var createdAtStr string
        
        err := rows.Scan(&task.ID, &task.Title, &task.Completed, &createdAtStr)
        if err != nil {
            log.Printf("❌ Ошибка сканирования строки: %v", err)
            continue
        }
        
        // Парсим время создания
        task.CreatedAt = parseTime(createdAtStr)
        tasks = append(tasks, task)
    }
    
    // Получаем общее количество для пагинации
    var total int
    countQuery := "SELECT COUNT(*) FROM tasks"
    if len(conditions) > 0 {
        countQuery += " WHERE " + strings.Join(conditions, " AND ")
    }
    
    err = r.db.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)
    if err != nil {
        total = len(tasks)
    }
    
    return tasks, total, nil
}

func (r *TaskRepository) FindByID(id int) (*models.Task, error) {
    var task models.Task
    var createdAtStr string
    
    err := r.db.QueryRow(
        "SELECT id, title, completed, created_at FROM tasks WHERE id = ?",
        id,
    ).Scan(&task.ID, &task.Title, &task.Completed, &createdAtStr)
    
    if err == sql.ErrNoRows {
        return nil, nil // Задача не найдена
    }
    
    if err != nil {
        return nil, fmt.Errorf("ошибка получения задачи %d: %v", id, err)
    }
    
    task.CreatedAt = parseTime(createdAtStr)
    return &task, nil
}

func (r *TaskRepository) Create(task *models.TaskRequest) (int64, error) {
    result, err := r.db.Exec(
        "INSERT INTO tasks (title, completed) VALUES (?, ?)",
        task.Title, task.Completed,
    )
    
    if err != nil {
        return 0, fmt.Errorf("ошибка создания задачи: %v", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("ошибка получения ID созданной задачи: %v", err)
    }
    
    return id, nil
}

func (r *TaskRepository) Update(id int, task *models.TaskRequest) error {
    result, err := r.db.Exec(
        "UPDATE tasks SET title = ?, completed = ? WHERE id = ?",
        task.Title, task.Completed, id,
    )
    
    if err != nil {
        return fmt.Errorf("ошибка обновления задачи %d: %v", id, err)
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("задача с ID %d не найдена", id)
    }
    
    return nil
}

func (r *TaskRepository) Delete(id int) error {
    result, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
    if err != nil {
        return fmt.Errorf("ошибка удаления задачи %d: %v", id, err)
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("задача с ID %d не найдена", id)
    }
    
    return nil
}

func (r *TaskRepository) GetStats() (*models.TaskStats, error) {
    var stats models.TaskStats
    
    // Получаем общее количество задач
    err := r.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&stats.Total)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения статистики: %v", err)
    }
    
    // Получаем количество выполненных задач
    err = r.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE completed = 1").Scan(&stats.Completed)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения статистики: %v", err)
    }
    
    stats.Active = stats.Total - stats.Completed
    return &stats, nil
}

func (r *TaskRepository) CheckHealth() error {
    return r.db.Ping()
}

func (r *TaskRepository) Exists(id int) (bool, error) {
    var exists bool
    err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", id).Scan(&exists)
    return exists, err
}

func parseTime(timeStr string) time.Time {
    if timeStr == "" {
        return time.Now()
    }
    
    if t, err := time.Parse("2006-01-02 15:04:05", timeStr); err == nil {
        return t
    }
    
    return time.Now()
}