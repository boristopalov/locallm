package db

import (
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type MilvusClient = struct {
	client.Client
	context.Context
}

var milvusClient MilvusClient = MilvusClient{
	nil,
	context.Background(),
}

func Client() (MilvusClient, error) {
	if milvusClient.Client != nil {
		return milvusClient, nil
	}
	client, err := client.NewClient(context.Background(), client.Config{
		Address: "localhost:19530", // TODO change
	})
	milvusClient.Client = client
	if err != nil {
		return milvusClient, err
	}
	return milvusClient, nil
}
