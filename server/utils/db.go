package utils

import (
	"context"
	"fmt"

	"github.com/boristopalov/localsearch/types"
	"github.com/boristopalov/localsearch/vars"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusClient = struct {
	client.Client
	context.Context
}

var milvusClient MilvusClient = MilvusClient{
	nil,
	context.Background(),
}

func GetDbClient() (MilvusClient, error) {
	if milvusClient.Client != nil {
		return milvusClient, nil
	}
	client, err := client.NewClient(context.Background(), client.Config{
		Address: "localhost:19530", // TODO change
	})
	client.UsingDatabase(milvusClient.Context, vars.DB_NAME)
	milvusClient.Client = client
	if err != nil {
		return milvusClient, err
	}
	return milvusClient, nil
}

func SaveTextEmbeddingsToVectorDB(we []types.TextEmbedding) error {
	c, err := GetDbClient()
	if err != nil {
		return err
	}

	err = c.LoadCollection(c.Context, vars.COLLECTION_NAME, false)
	if err != nil {
		return err
	}

	fmt.Println("collection loaded")
	hashes := make([]string, 0, len(we))
	names := make([]string, 0, len(we))
	texts := make([]string, 0, len(we))
	urls := make([]string, 0, len(we))
	embeddings := make([][]float32, 0, len(we))
	for i, item := range we {
		hashes = append(hashes, item.TextHash)
		names = append(names, item.Name)
		texts = append(texts, item.Text)
		urls = append(urls, item.Url)
		embeddings = append(embeddings, we[i].Embedding[:])
	}
	fmt.Println("creating fields to upsert")

	// bulk upsert since Milvus allows duplicates for some reason
	nameColumn := entity.NewColumnVarChar("Name", names)
	urlColumn := entity.NewColumnVarChar("Url", urls)
	textColumn := entity.NewColumnVarChar("Text", texts)
	hashColumn := entity.NewColumnVarChar("TextHash", hashes)
	vectorColumn := entity.NewColumnFloatVector("Embedding", vars.EMBEDDING_DIMS, embeddings)
	_, err = c.Upsert(c.Context, vars.COLLECTION_NAME, "", nameColumn, textColumn, hashColumn, urlColumn, vectorColumn)
	if err != nil {
		fmt.Println("failed to upsert data:", err.Error())
	}
	fmt.Println("upserted data")
	err = c.Flush(c.Context, vars.COLLECTION_NAME, false)
	if err != nil {
		return err
	}
	fmt.Println("flushed")
	err = c.ReleaseCollection(c.Context, vars.COLLECTION_NAME)
	if err != nil {
		return err
	}
	fmt.Println("released")
	return nil
}

func MaybeCreateCollection() error {
	c, err := GetDbClient()
	if err != nil {
		fmt.Println("error getting DB client: ", err.Error())
		return err
	}
	collExists, err := c.HasCollection(c.Context, vars.COLLECTION_NAME)
	if err != nil {
		return err
	}
	if !collExists {
		fmt.Println("cannot find collection... creating it")
		schema := entity.NewSchema().WithName(vars.COLLECTION_NAME).WithDescription("localsearch collection").
			WithField(entity.NewField().WithName("Name").WithMaxLength(65535).WithDataType(entity.FieldTypeVarChar)).
			WithField(entity.NewField().WithName("Text").WithMaxLength(65535).WithDataType(entity.FieldTypeVarChar)).
			WithField(entity.NewField().WithName("TextHash").WithMaxLength(65535).WithDataType(entity.FieldTypeVarChar).WithIsPrimaryKey(true)).
			WithField(entity.NewField().WithName("Url").WithMaxLength(65535).WithDataType(entity.FieldTypeVarChar)).
			WithField(entity.NewField().WithName("Embedding").WithDataType(entity.FieldTypeFloatVector).WithDim(vars.EMBEDDING_DIMS))
		err := c.CreateCollection(c.Context, schema, 2)
		idx, _ := entity.NewIndexFlat(entity.COSINE)
		// i := entity.NewScalarIndex()
		// c.CreateIndex(c.Context, vars.COLLECTION_NAME, "Name")
		c.CreateIndex(c.Context, vars.COLLECTION_NAME, "Embedding", idx, false)
		if err != nil {
			return err
		}
	}
	return nil
}
