package config

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDatabase() error {
    dbPath := AppConfig.DBPath
    log.Printf("📁 Используется файл БД: %s", dbPath)
    
    // Открываем подключение к SQLite базе данных
    db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_sync=1&_timeout=5000")
    if err != nil {
        return fmt.Errorf("ошибка подключения к базе данных %s: %v", dbPath, err)
    }
    
    // Устанавливаем максимальное количество открытых соединений
    db.SetMaxOpenConns(1)
    db.SetMaxIdleConns(1)
    
    // Проверяем соединение
    if err = db.Ping(); err != nil {
        return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
    }
    
    DB = db
    
    // Создаем таблицы
    if err := createTables(); err != nil {
        return err
    }
    
    fmt.Println("✅ База данных успешно инициализирована")
    return nil
}

func createTables() error {
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        completed BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
    
    createIndexSQL := `
    CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(completed);
    CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);`
    
    if _, err := DB.Exec(createTableSQL); err != nil {
        return fmt.Errorf("ошибка создания таблицы tasks: %v", err)
    }
    
    if _, err := DB.Exec(createIndexSQL); err != nil {
        log.Println("⚠️ Предупреждение: не удалось создать индексы:", err)
    }
    
    return nil
}