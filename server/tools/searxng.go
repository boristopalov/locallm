package tools

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/boristopalov/localsearch/utils"
)

type SearchTool interface {
	Execute(q string) (string, error)
}

type WebSearchTool struct {
	Action      string
	Description string
}

type SearchResult struct {
	URL           string   `json:"url"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	ImgSrc        string   `json:"img_src"`
	Engine        string   `json:"engine"`
	ParsedURL     []string `json:"parsed_url"`
	Template      string   `json:"template"`
	Engines       []string `json:"engines"`
	Positions     []int    `json:"positions"`
	PublishedDate string   `json:"publishedDate"`
	Score         float32  `json:"score"`
	Category      string   `json:"category"`
}

type SearchResponse struct {
	Query               string         `json:"query"`
	NumberOfResults     int            `json:"number_of_results"`
	Results             []SearchResult `json:"results"`
	Suggestions         []string       `json:"suggestions"`
	UnresponsiveEngines []interface{}  `json:"unresponsive_engines"`
}

// searches the web and returns the top results
func (t WebSearchTool) Execute(q string) (string, error) {
	httpClient := utils.HttpClient()
	// query, _ := json.Marshal(q)
	fmt.Println("search URI:", utils.SEARXNG_SEARCH_URI)
	resp, err := httpClient.Get(utils.SEARXNG_SEARCH_URI + fmt.Sprintf("?q=%s&format=json", url.QueryEscape(q)))
	// resp, err := httpClient.Post(utils.SEARXNG_SEARCH_URI, "application/json", bytes.NewReader(query))
	if err != nil {
		fmt.Println("error searching with searxng:", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	var searchResponse SearchResponse
	if err != nil {
		fmt.Println("error reading bytes")
		return "", err
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&searchResponse)
	if err != nil {
		fmt.Println("error unmarshaling json:", err.Error())
		return "", err
	}
	return GetTopResults(searchResponse), nil
}

// returns an array of the top results
// also saves top results to vector B
func GetTopResults(res SearchResponse) string {
	numTopResults := 10

	var topThreeResults strings.Builder
	splitOptions := utils.SplitOptions{
		MinLength:  10,
		MaxLength:  1000,
		Overlap:    20,
		Splitter:   "",
		Delimiters: "",
	}
	for i := 0; i < numTopResults; i++ {
		websiteUrl := res.Results[i].URL
		if strings.HasSuffix(websiteUrl, ".pdf") {
			continue
		}
		html, err := utils.GetWebsiteHTML(websiteUrl)
		if err != nil {
			continue
		}
		htmlTextOnly := utils.ExtractText(html) // we only need the text from the html

		fullTextChunk := utils.WebsiteData{
			URL:   websiteUrl,
			Title: res.Results[1].Title,
			Text:  htmlTextOnly,
		}
		fullTextChunkJSON, err := json.Marshal(fullTextChunk)
		if err != nil {
			fmt.Println("error marshaling json to string", err.Error())
			return ""
		}

		// only the top 3 results will get returned to the LLM
		// the rest just get saved in the vector DB
		if i < 3 {
			topThreeResults.WriteString(string(fullTextChunkJSON) + "\n-------------\n")
		}

		chunks := utils.Chunk(htmlTextOnly, splitOptions) // split the website text into chunks
		for _, c := range chunks {
			websiteChunk := utils.WebsiteData{
				URL:   websiteUrl,
				Title: res.Results[i].Title,
				Text:  c,
			}
			embedding, err := utils.EmbedWebsiteChunk(websiteChunk) // create embedding
			if err != nil {
				continue
			}
			SaveTextEmbeddingToVectorDB(embedding)
		}
	}
	fmt.Println("Top three results returned to the model:", topThreeResults.String())
	return topThreeResults.String()
}

var WebSearchTool_ = WebSearchTool{
	Action: "WebSearch",
	Description: `Useful for searching the internet. 
				You must use this tool if you're not 100% certain of the answer. 
				The top 10 results will be added to the vector db. 
				The top 3 results are also getting returned to you directly as JSON. 
				For more search queries through the same websites, use the VectorDB tool`,
}
