package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
	_ "github.com/lib/pq"
)

type DBConfing struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

func Init(config DBConfing) (*sql.DB, *sqlcgen.Queries, error) {
	CONNECT := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
	)
	db, err := sql.Open("postgres", CONNECT)
	if err != nil {
		log.Fatal("Failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("faild to ping detabase: %w", err)
	}

	queries := sqlcgen.New(db)
	return db, queries, nil
}
