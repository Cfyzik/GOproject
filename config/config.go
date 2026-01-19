package config

import (
    "log"
    "os"
    
    "github.com/joho/godotenv"
)

type Config struct {
    DBPath  string
    Port    string
    GinMode string
}

var AppConfig Config

func Init() {
    // Загружаем .env файл
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ .env файл не найден, используются значения по умолчанию")
    }
    
    // Читаем конфигурацию из переменных окружения
    AppConfig = Config{
        DBPath:  getEnv("DB_PATH", "./tasks.db"),
        Port:    getEnv("PORT", "8080"),
        GinMode: getEnv("GIN_MODE", "debug"),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}