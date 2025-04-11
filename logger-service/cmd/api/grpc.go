package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	log.Printf("Received log entry: %+v", input)

	// write the log
	logEntry := data.LogEntry{
		Name: input.GetName(),
		Data: input.GetData(),
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		log.Printf("Error inserting logEntry: %v", err)
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// return response
	res := &logs.LogResponse{Result: "logged!"}

	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC server listening on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
