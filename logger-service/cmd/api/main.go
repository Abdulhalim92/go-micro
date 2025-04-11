package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50051"
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

	// close connection
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			log.Panicf("Error disconnecting from mongo: %v", err)
		}
	}()

	app := Config{
		Models: data.New(mongoClient),
	}

	// register RPC server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panicf("Error registering RPC server: %v", err)
	}
	go app.rpcListen()

	go app.gRPCListen()

	// start the web server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Listening on port %s", webPort)

	err = srv.ListenAndServe()
	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port", rpcPort)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		rpcConn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting RPC connection: %v", err)
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect to mongo
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Printf("Error connecting to mongo: %v", err)
		return nil, err
	}

	log.Println("Connected to mongo")

	return client, nil
}
