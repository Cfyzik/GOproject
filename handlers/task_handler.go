package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "strings"
    
    "github.com/gin-gonic/gin"
    
    "GOproject/models"
    "GOproject/services"
)

type TaskHandler struct {
    service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
    return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
    // Получаем параметры фильтрации
    status := c.Query("status")
    search := c.Query("search")
    
    limit := getIntQuery(c, "limit", 100, 1, 1000)
    offset := getIntQuery(c, "offset", 0, 0, 10000)
    
    result, err := h.service.GetAllTasks(status, search, limit, offset)
    if err != nil {
        log.Printf("❌ Ошибка получения задач: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
        return
    }
    
    c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
        return
    }
    
    task, err := h.service.GetTaskByID(id)
    if err != nil {
        if _, ok := err.(*services.NotFoundError); ok {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
        return
    }
    
    if task == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
        return
    }
    
    c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
    var taskReq models.TaskRequest
    
    // Используем json.Decoder для лучшей обработки ошибок
    decoder := json.NewDecoder(c.Request.Body)
    decoder.DisallowUnknownFields()
    
    if err := decoder.Decode(&taskReq); err != nil {
        var errMsg string
        if strings.Contains(err.Error(), "unknown field") {
            errMsg = "Недопустимое поле в запросе"
        } else {
            errMsg = "Неверный формат JSON"
        }
        
        c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
        return
    }
    defer c.Request.Body.Close()
    
    task, err := h.service.CreateTask(&taskReq)
    if err != nil {
        if validationErr, ok := err.(*services.ValidationError); ok {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }
        
        log.Printf("❌ Ошибка создания задачи: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать задачу"})
        return
    }
    
    // Логируем успешное создание
    log.Printf("✅ Создана новая задача: ID=%d, Title='%s'", task.ID, task.Title)
    
    c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
        return
    }
    
    var taskReq models.TaskRequest
    
    decoder := json.NewDecoder(c.Request.Body)
    decoder.DisallowUnknownFields()
    
    if err := decoder.Decode(&taskReq); err != nil {
        var errMsg string
        if strings.Contains(err.Error(), "unknown field") {
            errMsg = "Недопустимое поле в запросе"
        } else {
            errMsg = "Неверный формат JSON"
        }
        
        c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
        return
    }
    defer c.Request.Body.Close()
    
    task, err := h.service.UpdateTask(id, &taskReq)
    if err != nil {
        if validationErr, ok := err.(*services.ValidationError); ok {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }
        
        if _, ok := err.(*services.NotFoundError); ok {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        
        log.Printf("❌ Ошибка обновления задачи %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить задачу"})
        return
    }
    
    // Логируем успешное обновление
    log.Printf("🔄 Обновлена задача: ID=%d, Title='%s', Completed=%v", task.ID, task.Title, task.Completed)
    
    c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
        return
    }
    
    err = h.service.DeleteTask(id)
    if err != nil {
        if _, ok := err.(*services.NotFoundError); ok {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        
        log.Printf("❌ Ошибка удаления задачи %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить задачу"})
        return
    }
    
    // Логируем успешное удаление
    log.Printf("🗑️ Удалена задача: ID=%d", id)
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Задача успешно удалена",
        "id":      id,
    })
}

func getIntQuery(c *gin.Context, key string, defaultValue, min, max int) int {
    str := c.Query(key)
    if str == "" {
        return defaultValue
    }
    
    val, err := strconv.Atoi(str)
    if err != nil || val < min || val > max {
        return defaultValue
    }
    
    return val
}