package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	var err error
	client, err = connectToMonge()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}
	//Register the rpc server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		go app.rpcListen()
	}
	go app.grpcListen()
	app.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("starting rpc server on port: ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMonge() (*mongo.Client, error) {

	clientOpt := options.Client().ApplyURI(mongoURL)
	clientOpt.SetAuth(
		options.Credential{
			Username: "admin",
			Password: "password",
		},
	)
	c, err := mongo.Connect(clientOpt)
	if err != nil {
		log.Println("Error connecting, ", err)
		return nil, err
	}
	return c, nil

}
