// Файл: internal/db/db.go
// Назначение: установка подключения к Postgres через sqlx и применение SQL-миграций из internal/db/migrations.
package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect открывает пул соединений с Postgres по заданной строке DSN.
func Connect(dsn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	return db
}

// Migrate применяет все пендинговые SQL-миграции из директории internal/db/migrations.
func Migrate(dsn string) {
	m, err := migrate.New("file://internal/db/migrations", dsn)
	if err != nil {
		log.Fatalf("migrate init: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up: %v", err)
	}
	log.Println("migrations applied")
}
