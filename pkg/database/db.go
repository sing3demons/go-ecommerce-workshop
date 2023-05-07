package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sing3demons/shop/config"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	// Connect
	db, err := sqlx.Connect("pgx", cfg.Dsn())
	if err != nil {
		log.Fatalf("connect to db failed: %v\n", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())

	return db
}

func ConnectDB(cfg config.IDbConfig) (*sql.DB, error) {
	sql, err := sql.Open("postgres", cfg.Dsn())
	if err != nil {
		log.Printf("open db failed: %v", err)
	}

	sql.SetMaxOpenConns(cfg.MaxOpenConns())

	if err := sql.Ping(); err != nil {
		log.Printf("ping db failed: %v", err)
	}

	return sql, nil
}
