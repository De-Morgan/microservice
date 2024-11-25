package main

import (
	"authentication/data"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const webPort = "80"

type Config struct {
	DB     *pgxpool.Pool
	Models data.Models
}

func main() {

	log.Println("starting authentication service")
	conn := connectToDB()
	if conn == nil {
		log.Println("Can't connect to postgres")
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

func openDB(dsn string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil

}

func connectToDB() *pgxpool.Pool {
	dsn := os.Getenv("DSN")
	var count int
	for {
		cnt, err := openDB(dsn)
		if err != nil {
			log.Println("database not ready:")
			count++
		} else {
			log.Println("connected to postgress")
			return cnt
		}
		if count > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off 2 sec")
		time.Sleep(2 * time.Second)
		continue
	}

}
