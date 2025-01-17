package main

import (
	"context"
	"log"
	"log-service/data"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, res *string) error {

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.Background(), data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
		ID:   primitive.NewObjectID().Hex(),
	})

	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	*res = "Processed payload: " + payload.Name

	return nil

}
