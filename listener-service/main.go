package main

import (
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to rabbitmq
	conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create a consumer
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume messages
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	// @NOTE: rabbitmq takes time to start up, so we need to wait for it
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("RabbitMQ not yet ready...")
			counts++
			if counts > 5 {
				log.Println(err)
				return nil, err
			}
			backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
			log.Println("backing off...")
			time.Sleep(backOff)
			continue
		}
		log.Println("Connected to RabbitMQ!")
		connection = c
		break
	}

	return connection, nil
}
