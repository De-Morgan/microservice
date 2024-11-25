package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// Write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := data.Insert(logEntry)
	if err != nil {
		res := logs.LogResponse{Result: "failed"}
		return &res, err
	}

	res := logs.LogResponse{Result: "logged"}

	return &res, nil
}

func (c *Config) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()

	logs.RegisterLogServiceServer(srv, &LogServer{})
	log.Println("grpc server started on: ", grpcPort)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
