package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	rabitConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer rabitConn.Close()

	log.Println("connected to rabit")

	app := Config{
		Rabbit: rabitConn,
	}

	log.Printf("starting broker on port %q\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {

	var count int64
	backOff := 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabitmq")
		if err != nil {
			fmt.Println("rabitmq not ready:", err)
			count++
		} else {
			connection = c
			break
		}

		if count > 5 {
			fmt.Println("rabitmq can't connect:", err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("backing off: ", backOff)
		time.Sleep(backOff)
	}

	return connection, nil

}
