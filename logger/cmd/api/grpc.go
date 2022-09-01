package main

import (
	"context"
	"logger/data"
	"logger/logs"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(context context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	if err := l.Models.LogEntry.Insert(logEntry); err != nil {
		return &logs.LogResponse{Result: "Failed to write logs"}, err
	}

	return &logs.LogResponse{Result: "Success writing logs"}, nil
}
