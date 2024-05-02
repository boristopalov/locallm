package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boristopalov/localsearch/types"
	"github.com/boristopalov/localsearch/vars"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

/*
// Text Embedding helper functions below
*/
func GetTextEmbedding(text string) ([]float32, error) {
	var data = struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}{
		Model:  vars.EMBEDDINGS_MODEL,
		Prompt: text,
	}
	reqData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(vars.OLLAMA_EMBED_URI, "text/html", bytes.NewReader(reqData))
	if err != nil {
		fmt.Println("error sending request to embed endpoint:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	var respBody = struct {
		Embedding []float32 `json:"embedding"`
	}{
		Embedding: make([]float32, vars.EMBEDDING_DIMS),
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respBody)
	if err != nil {
		fmt.Println("error unmarshalling JSON", err.Error())
		return nil, err
	}
	return respBody.Embedding, nil
}

func EmbedWebsiteChunk(chunk types.WebsiteData) (types.TextEmbedding, error) {
	textEmbedding, err := GetTextEmbedding(chunk.Text)
	if err != nil {
		fmt.Println("error embedding text", err.Error())
		return types.TextEmbedding{}, err
	}
	var embedStruct = types.TextEmbedding{
		Text:      chunk.Text,
		TextHash:  GetMD5Hash(chunk.Text),
		Embedding: textEmbedding,
		Name:      chunk.Title,
		Url:       chunk.URL,
	}
	return embedStruct, nil
}

func FormatSSEMessage(event types.DataStreamEvent) string {
	e := fmt.Sprintf("event: %s\n", event.EventName)
	text := fmt.Sprintf("data: %s\n\n", base64.StdEncoding.EncodeToString([]byte(event.Data)))
	fmt.Println(e + text)
	return e + text
}
