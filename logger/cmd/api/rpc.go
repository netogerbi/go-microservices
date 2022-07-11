package main

import (
	"context"
	"log"
	"logger/data"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(p RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")

	le := data.LogEntry{
		Name:      p.Name,
		Data:      p.Data,
		CreatedAt: time.Now(),
	}

	if _, err := collection.InsertOne(context.TODO(), le); err != nil {
		log.Println(err)
		return err
	}

	*resp = "Processed payload via RPC: " + le.Name

	return nil
}
