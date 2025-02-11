package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	webPort  = "4000"
	rpcPort  = "4001"
	grpcPort = "4002"
	mongoURL = "mongodb://mongo:27017"
	// @NOTE: Running with docker use mongo:27017
	// mongo is the service name in the docker-compose.yaml file
	// 27017 is the container port for mongo
	// Running locally use localhost:27018
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close mongo connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// register the RPC server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	// start the server
	log.Println("Starting logger service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create a client options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return client, nil
}

func (app *Config) rpcListen() {
	log.Println("Starting RPC server on port", rpcPort)

	rpcListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Println(err)
		return
	}
	defer rpcListener.Close()

	for {
		rpcConn, err := rpcListener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
