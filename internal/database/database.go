package database

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"subscription-service/internal/config"
)

// Connect создает подключение к PostgreSQL и выполняет миграции
func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	log.Println("Подключение к базе данных установлено")

	// Выполнение версионных миграций
	migrationPath := filepath.Join("migrations")
	absPath, err := filepath.Abs(migrationPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить абсолютный путь миграций: %w", err)
	}
	m, err := migrate.New("file://"+absPath, "postgres://"+cfg.User+":"+cfg.Password+"@"+cfg.Host+":"+cfg.Port+"/"+cfg.Name+"?sslmode="+cfg.SSLMode)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации миграций: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("ошибка применения миграций: %w", err)
	}

	log.Println("Миграции базы данных выполнены")
	return db, nil
}
