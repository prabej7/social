package main

import (
	"log"

	idb "github.com/prabej7/social/internal/db"
	"github.com/prabej7/social/internal/env"
	"github.com/prabej7/social/internal/store"
)

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func main() {
	dbConfig := dbConfig{
		addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/social?sslmode=disable"),
		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15min"),
	}

	db, err := idb.New(
		dbConfig.addr,
		dbConfig.maxOpenConns,
		dbConfig.maxIdleConns,
		dbConfig.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection pool established")
	storage := store.NewStorage(db)

	idb.Seed(&storage)
}
