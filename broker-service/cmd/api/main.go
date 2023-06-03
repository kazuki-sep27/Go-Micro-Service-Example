package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct{
	Rabbit *amqp.Connection
}

func main() {
	//try to connect to rabbitmq
	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
		
	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting API server on port %s", webPort)

	// define http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start http server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func connectToRabbitMQ() (*amqp.Connection, error) {
	// connect to rabbitmq
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until we have a connection
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		}else{
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}
		
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Backing off for %v", backOff)
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}