package types

type ClientQuery struct {
	Prompt    string `json:"prompt"`
	Model     string `json:"model"`
	KeepAlive string `json:"keep_alive"`
}

type OllamaQuery struct {
	Prompt    string `json:"prompt"`
	Template  string `json:"template"`
	Model     string `json:"model"`
	KeepAlive string `json:"keep_alive"`
	Context   []int  `json:"context"`
	System    string `json:"system"`
}

type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	Context   []int  `json:"context"`
}

type Tool struct {
	Action      string
	Description string
}

type TextEmbedding struct {
	Text      string
	TextHash  string
	Name      string
	Url       string
	Embedding []float32
}

type WebsiteData struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DataStreamEvent struct {
	EventName string
	Data      string
}

type WebSearchResult struct {
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

type WebSearchResponse struct {
	Query               string            `json:"query"`
	NumberOfResults     int               `json:"number_of_results"`
	Results             []WebSearchResult `json:"results"`
	Suggestions         []string          `json:"suggestions"`
	UnresponsiveEngines []interface{}     `json:"unresponsive_engines"`
}

type VectorSearchResult struct {
	Text         string `json:"text"`
	WebsiteTitle string `json:"websiteTitle"`
	Url          string `json:"url"`
}

type ToolResponse struct {
	Data   string `json:"data"`
	Source string `json:"source"`
}
