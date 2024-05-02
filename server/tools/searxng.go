package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/boristopalov/localsearch/types"
	"github.com/boristopalov/localsearch/utils"
	"github.com/boristopalov/localsearch/vars"
)

type WebSearchToolInterface interface {
	Execute(q string) (string, error)
}
type WebSearchToolType struct{}

var WebSearchTool = WebSearchToolType{}

// searches the web, saves top (10) results in vector DB
// then searches the vector DB using the search string and returns top (3) results
func (t WebSearchToolType) Execute(q string) ([]types.WebsiteData, error) {
	// query, _ := json.Marshal(q)
	fmt.Println("search URI:", vars.SEARXNG_SEARCH_URI)
	resp, err := http.Get(vars.SEARXNG_SEARCH_URI + fmt.Sprintf("?q=%s&format=json", url.QueryEscape(q)))
	// resp, err := httpClient.Post(utils.SEARXNG_SEARCH_URI, "application/json", bytes.NewReader(query))
	if err != nil {
		fmt.Println("error searching with searxng:", err.Error())
		return nil, err
	}
	fmt.Println("search query ran successfully")
	defer resp.Body.Close()
	var searchResponse types.WebSearchResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&searchResponse)
	if err != nil {
		fmt.Println("error unmarshaling json:", err.Error())
		return nil, err
	}

	numTopResults := 10
	if len(searchResponse.Results) < numTopResults {
		numTopResults = len(searchResponse.Results)
	}

	var ret []types.WebsiteData
	for i := 0; i < numTopResults; i++ {
		res := types.WebsiteData{
			URL:   searchResponse.Results[i].URL,
			Title: searchResponse.Results[i].Title,
			Text:  searchResponse.Results[i].Content,
		}
		ret = append(ret, res)
	}
	go SaveTopResults(searchResponse)
	return ret, nil

	// if err != nil {
	// 	return "", err
	// }
	// fmt.Println("running similarity search on DB for string ", q)

	// topResults, err := QueryVectorDBTool.Execute(q)
	// if err != nil {
	// 	return "", fmt.Errorf("error in web search tool from vector db search: %s", err.Error())
	// }
	// fmt.Println("top results from web search + vectordb query: ", topResults)
	// return topResults, nil
}

// returns an array of the top results
// also saves top results to vector B
func SaveTopResults(res types.WebSearchResponse) error {
	fmt.Println("num results:", len(res.Results))
	fmt.Println("saving top results to vector DB...")
	numTopResults := 10
	if len(res.Results) < numTopResults {
		numTopResults = len(res.Results)
	}

	for i := 0; i < numTopResults; i++ {
		websiteUrl := res.Results[i].URL
		if strings.HasSuffix(websiteUrl, ".pdf") {
			continue
		}
		html, err := GetWebsiteHTML(websiteUrl)
		if err != nil {
			continue
		}
		htmlTextOnly := utils.ExtractText(html) // we only need the text from the html
		if len(htmlTextOnly) == 0 {
			fmt.Println("no text to extract!")
			continue
		}

		textSplitter := utils.DefaultRecursiveTextSplitter()
		chunks, err := textSplitter.SplitText(htmlTextOnly)
		if err != nil {
			continue
		}
		toInsert := make([]types.TextEmbedding, 0, len(chunks))
		fmt.Println("embedding chunks")
		for _, c := range chunks {
			websiteChunk := types.WebsiteData{
				URL:   websiteUrl,
				Title: res.Results[i].Title,
				Text:  c,
			}
			embedding, err := utils.EmbedWebsiteChunk(websiteChunk) // create embedding
			if err != nil {
				continue
			} else {
				toInsert = append(toInsert, embedding)
			}
		}
		fmt.Println("saving website to db")
		err = utils.SaveTextEmbeddingsToVectorDB(toInsert)
		if err != nil {
			fmt.Printf("failed to save embedding to vector DB: %s", err.Error())
			continue
		}
	}
	return nil
}

func GetWebsiteHTML(url string) (string, error) {
	fmt.Println("grabbing website HTML")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error accessing URL", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading HTML", err.Error())
		return "", err
	}
	return string(body), nil
}
