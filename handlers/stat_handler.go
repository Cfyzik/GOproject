package handlers

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    
    "GOproject/services"
)

type StatHandler struct {
    service *services.TaskService
}

func NewStatHandler(service *services.TaskService) *StatHandler {
    return &StatHandler{service: service}
}

func (h *StatHandler) HealthCheck(c *gin.Context) {
    // Проверяем соединение с БД
    if err := h.service.HealthCheck(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status":  "error",
            "message": "Database connection failed",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status":    "ok",
        "database":  "connected",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

func (h *StatHandler) GetStats(c *gin.Context) {
    stats, err := h.service.GetStats()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения статистики"})
        return
    }
    
    c.JSON(http.StatusOK, stats)
}