package utils

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
	Name      string
	Url       string
	Embedding []float32
}

type WebsiteData struct {
	URL   string
	Title string
	Text  string
}
