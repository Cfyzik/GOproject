package main

import (
    "fmt"
    "log"
    
    "GOproject/config"
    "GOproject/routes"
)

func main() {
    // Инициализация конфигурации
    config.Init()
    
    // Инициализация базы данных
    if err := config.InitDatabase(); err != nil {
        log.Fatal("❌ Ошибка инициализации базы данных:", err)
    }
    defer config.DB.Close()
    
    // Настройка роутера
    router := routes.SetupRouter()
    
    // Запуск сервера
    addr := ":" + config.AppConfig.Port
    fmt.Println("🚀 Сервер запущен!")
    fmt.Printf("📊 Адрес: http://localhost:%s\n", config.AppConfig.Port)
    fmt.Printf("🔧 API доступен по: http://localhost:%s/api/tasks\n", config.AppConfig.Port)
    fmt.Printf("📁 Файл БД: %s\n", config.AppConfig.DBPath)
    fmt.Println("⏳ Ожидание запросов...")
    
    // Запускаем сервер
    if err := router.Run(addr); err != nil {
        log.Fatal("Ошибка запуска сервера:", err)
    }
}