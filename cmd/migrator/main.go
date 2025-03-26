package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Подключаем PostgreSQL драйвер
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Источник миграций
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	// Принимаем аргументы командной строки
	flag.StringVar(&storagePath, "storage-path", "", "path to storage (e.g. PostgreSQL connection string)")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations (e.g. path to migrations directory)")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	// Проверка наличия обязательных параметров
	if storagePath == "" {
		panic("storage-path is required") 
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	// Создаем мигратор
	m, err := migrate.New(
		"file://"+migrationsPath, // Источник миграций (путь к миграциям)
		storagePath,
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create migrator: %v", err))
	}

	// Применяем миграции
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			// Если нет новых миграций, выводим сообщение
			fmt.Println("no migrations to apply")
			return
		}
		panic(fmt.Sprintf("Failed to apply migrations: %v", err))
	}

	// Если миграции были применены успешно
	fmt.Println("Migrations applied successfully")
}
