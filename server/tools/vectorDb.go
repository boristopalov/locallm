package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/boristopalov/localsearch/db"
	"github.com/boristopalov/localsearch/utils"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type VectorDBTool interface {
	Execute() ([]client.SearchResult, error)
}
type QueryVectorDBTool struct {
	Action      string
	Description string
}

func (t QueryVectorDBTool) Execute(q []float32) (string, error) {
	c, err := db.Client()
	if err != nil {
		fmt.Println("error getting db client", err.Error())
		return "", err
	}

	// do i need to specify params?
	sp, _ := entity.NewIndexFlatSearchParam()
	var topThreeResults strings.Builder

	// get 3 closest results
	searchResult, err := c.Client.Search(
		c.Context,
		utils.COLLECTION_NAME,
		[]string{},
		"",
		[]string{"Text"},
		[]entity.Vector{entity.FloatVector(q)},
		"Embedding",
		entity.COSINE,
		3,
		sp,
	)
	if err != nil {
		fmt.Println("failed to search collection", err.Error())
		return "", err
	}
	err = c.Client.ReleaseCollection(c.Context, utils.COLLECTION_NAME)
	if err != nil {
		fmt.Println("failed to release collection", err.Error())
		return "", err
	}
	for _, sr := range searchResult {
		out, err := json.Marshal(sr)
		if err != nil {
			fmt.Println("error marshaling JSON", err.Error())
			continue
		}
		topThreeResults.WriteString(string(out) + "\n-------------\n")
	}
	return topThreeResults.String(), nil
}

func SaveTextEmbeddingToVectorDB(we utils.TextEmbedding) error {
	c, err := db.Client()
	if err != nil {
		fmt.Println("error getting DB client", err.Error())
		return err
	}
	collExists, err := c.HasCollection(c.Context, utils.COLLECTION_NAME)
	if err != nil {
		log.Fatal("failed to check collection exists:", err.Error())
	}
	if !collExists {
		fmt.Println("cannot find collection... creating it")
		schema := entity.NewSchema().WithName(utils.COLLECTION_NAME).WithDescription("localsearch collection").
			WithField(entity.NewField().WithName("Name").WithMaxLength(65535).WithDataType(entity.FieldTypeVarChar)).
			WithField(entity.NewField().WithName("Url").WithMaxLength(65535).WithDataType(entity.FieldTypeVarChar).WithIsPrimaryKey(true)).
			WithField(entity.NewField().WithName("Embedding").WithDataType(entity.FieldTypeFloatVector).WithDim(utils.EMBEDDING_DIMS))
		err := c.CreateCollection(c.Context, schema, 2)
		if err != nil {
			log.Fatal("failed to create collection:", err.Error())
		}
	}
	found, err := c.Get(c.Context, utils.COLLECTION_NAME, entity.NewColumnString("Url", []string{we.Url}))
	if err != nil {
		fmt.Println("error calling Get() on db", err.Error())
		return err
	}
	if found != nil {
		fmt.Println("website already saved in db")
		return nil
	}
	// TODO: why are these all arrays
	name := entity.NewColumnVarChar("Name", []string{we.Name})
	url := entity.NewColumnVarChar("Url", []string{we.Url})
	vector := entity.NewColumnFloatVector("Embedding", utils.EMBEDDING_DIMS, [][]float32{we.Embedding})
	_, err = c.Insert(c.Context, utils.COLLECTION_NAME, "", name, url, vector)
	if err != nil {
		log.Fatal("failed to insert data:", err.Error())
	}
	return nil
}

var QueryVectorDBTool_ = QueryVectorDBTool{
	Action: "QueryVectorDB",
	Description: `Useful for searching through added files and websites. 
	Search for keywords in the text, not whole questions. 
	Avoid relative words like "yesterday" and think about what could be in the text. 
	The input to this tool will be run against a vector db. The top results will be returned to you as JSON.`,
}
