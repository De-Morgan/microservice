package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	rabitConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer rabitConn.Close()

	log.Println("connected to rabit")

	consumer, err := event.NewConsumer(rabitConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})

	if err != nil {
		panic(err)
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
