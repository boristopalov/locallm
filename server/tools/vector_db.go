package tools

import (
	"encoding/json"
	"fmt"

	"github.com/boristopalov/localsearch/types"
	"github.com/boristopalov/localsearch/utils"
	"github.com/boristopalov/localsearch/vars"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type VectorDBTool interface {
	Execute(q string) (string, error)
}
type QueryVectorDBToolType struct{}

var QueryVectorDBTool = QueryVectorDBToolType{}

func (t QueryVectorDBToolType) Execute(q string) (string, error) {
	vectorQuery, err := utils.GetTextEmbedding(q)
	if err != nil {
		fmt.Println("error getting text embedding:", err.Error())
		return "", err
	}
	c, err := utils.GetDbClient()
	if err != nil {
		fmt.Println("error getting db client", err.Error())
		return "", err
	}
	err = c.LoadCollection(c.Context, vars.COLLECTION_NAME, false)
	if err != nil {
		fmt.Println("failed to load collection: ", err.Error())
	}

	// do i need to specify params?
	sp, _ := entity.NewIndexFlatSearchParam()

	fmt.Println("constructing vector search query...")
	// get 3 closest results
	topk := 10
	searchResult, err := c.Search(
		c.Context,
		vars.COLLECTION_NAME,
		[]string{},
		"",
		[]string{"Text", "Name", "Url"},
		[]entity.Vector{entity.FloatVector(vectorQuery)},
		"Embedding",
		entity.COSINE,
		topk,
		sp,
	)
	if err != nil {
		fmt.Println("failed to search collection: ", err.Error())
		return "", err
	}

	vsrs := make([]types.VectorSearchResult, 0, topk)
	var titles []string
	var texts []string
	var urls []string
	for _, sr := range searchResult {
		for _, field := range sr.Fields {
			s, ok := field.(*entity.ColumnVarChar)
			if ok {
				if field.Name() == "Name" {
					titles = s.Data()
				}
				if field.Name() == "Text" {
					texts = s.Data()
				}
				if field.Name() == "Url" {
					urls = s.Data()
				}
			}
		}
	}
	for i := 0; i < len(texts); i++ {
		vsr := types.VectorSearchResult{
			Text:         texts[i],
			WebsiteTitle: titles[i],
			Url:          urls[i],
		}
		vsrs = append(vsrs, vsr)
	}
	json, err := json.Marshal(vsrs)
	if err != nil {
		return "", nil
	}

	err = c.ReleaseCollection(c.Context, vars.COLLECTION_NAME)
	if err != nil {
		fmt.Println("failed to release collection: ", err.Error())
		return "", err
	}
	return string(json), nil
}
