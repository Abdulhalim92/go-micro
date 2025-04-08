package main

import (
	"authorization/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting the authentication service")

	// подключение к БД
	conn := connectToDB()
	if conn == nil {
		log.Panicln("Unable to connect to database")
	}

	// инициализация конфигурации
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// инициализация модели сервера
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// openDB - это функция, которая открывает соединение с базой данных
func openDB(dsn string) (*sql.DB, error) {
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

// connDB - это функция, которая устанавливает соединение с базой данных
func connectToDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet")
			counts++
		} else {
			log.Println("Connect to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Printf("Too many postgres connections: %v", err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
