package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"logger-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	return &logs.LogResponse{Result: "logged"}, nil
}

func (app *Config) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for grpc: %v", err)
	}
	defer listen.Close()

	grpcServer := grpc.NewServer()

	logs.RegisterLogServiceServer(grpcServer, &LogServer{Models: app.Models})

	log.Printf("gRPC server started on port %s", grpcPort)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to listen for grpc: %v", err)
	}
}
