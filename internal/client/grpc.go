package client

import (
	"github.com/AnandRaj2224/agentbox/internal/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Connect(target string) (api.ExecutionServiceClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return api.NewExecutionServiceClient(conn), nil
}
