package database

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

func NewPostgresDB(host, port, user, password, dbname string) (*sql.DB, error) {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    log.Println("Successfully connected to PostgreSQL")
    return db, nil
}