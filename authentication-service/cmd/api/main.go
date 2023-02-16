package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "80"

var countTimeout = 0

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("Starting authentication server on port %s", webPort)

	conn := connectDb()
	if conn == nil {
		log.Panic("Cant connect to postgres!")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDb() *sql.DB {
	dns := os.Getenv("DATABASE_DSN")

	for {
		conn, err := openDb(dns)
		if err != nil {
			log.Println("Postgres is not ready yet...")
		} else {
			log.Println("Connected to Postgres!")
			return conn
		}

		if countTimeout > 10 {
			log.Panic(err)
			return nil
		}

		fmt.Println("Waiting 2 seconds for DB to get ready.")
		countTimeout++
		time.Sleep(2 * time.Second)
		continue
	}
}
