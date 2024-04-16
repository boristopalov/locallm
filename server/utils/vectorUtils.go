package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func GetWebsiteHTML(url string) (string, error) {
	httpClient := HttpClient()
	resp, err := httpClient.Get(url)
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

func GetTextEmbedding(text string) ([]float32, error) {
	httpClient := HttpClient()
	var data = struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}{
		Model:  EMBEDDINGS_MODEL,
		Prompt: text,
	}
	reqData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Post(OLLAMA_EMBED_URI, "text/html", bytes.NewReader(reqData))
	if err != nil {
		fmt.Println("error sending request to embed endpoint:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	var respBody = struct {
		Embedding []float32 `json:"embedding"`
	}{
		Embedding: make([]float32, EMBEDDING_DIMS),
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respBody)
	if err != nil {
		fmt.Println("error unmarshalling JSON", err.Error())
		return nil, err
	}
	return respBody.Embedding, nil
}

func EmbedWebsiteChunk(chunk WebsiteData) (TextEmbedding, error) {
	textEmbedding, err := GetTextEmbedding(chunk.Text)
	if err != nil {
		fmt.Println("error embedding text", err.Error())
		return TextEmbedding{}, err
	}
	var embedStruct = struct {
		Name      string
		Url       string
		Embedding []float32
	}{
		Embedding: textEmbedding,
		Name:      chunk.Title,
		Url:       chunk.URL,
	}
	return embedStruct, nil
}
