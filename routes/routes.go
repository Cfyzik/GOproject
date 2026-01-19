package routes

import (
    "github.com/gin-gonic/gin"
    
    "GOproject/config"
    "GOproject/handlers"
    "GOproject/repository"
    "GOproject/services"
)

func SetupRouter() *gin.Engine {
    // Устанавливаем режим Gin
    if config.AppConfig.GinMode == "release" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    router := gin.Default()
    
    // Middleware
    router.Use(handlers.CORSMiddleware())
    
    // Статические файлы
    router.Static("/static", "./static")
    
    // Главная страница
    router.GET("/", func(c *gin.Context) {
        c.File("./static/index.html")
    })
    
    // Инициализируем зависимости
    taskRepo := repository.NewTaskRepository()
    taskService := services.NewTaskService(taskRepo)
    taskHandler := handlers.NewTaskHandler(taskService)
    statHandler := handlers.NewStatHandler(taskService)
    
    // API маршруты
    api := router.Group("/api")
    {
        // Задачи
        api.GET("/tasks", taskHandler.GetTasks)
        api.GET("/tasks/:id", taskHandler.GetTask)
        api.POST("/tasks", taskHandler.CreateTask)
        api.PUT("/tasks/:id", taskHandler.UpdateTask)
        api.DELETE("/tasks/:id", taskHandler.DeleteTask)
        
        // Статистика
        api.GET("/health", statHandler.HealthCheck)
        api.GET("/stats", statHandler.GetStats)
    }
    
    return router
}