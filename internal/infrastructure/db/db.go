package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
	"github.com/cockroachdb/errors"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

func Init(config DBConfig) (*sql.DB, *sqlcgen.Queries, error) {
	var err error
	var db *sql.DB
	CONNECT := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
	)

	// 接続できるまで一定回数リトライ
	count := 0
	db, err = sql.Open("postgres", CONNECT)
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 { // countが180になるまでリトライ
				fmt.Println("")
				return db, nil, errors.WithStack(errors.New("db Inin error. count over 180."))
			}
			db, err = sql.Open("postgres", CONNECT)
		}
	}

	queries := sqlcgen.New(db)
	return db, queries, nil
}
