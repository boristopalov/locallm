package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
 */
func PromptModel(query OllamaQuery, ch chan<- OllamaResponse) (string, error) {
	data, err := json.Marshal(query)
	if err != nil {
		fmt.Println("error marshalling data:", query)
		return "", err
	}
	req, err := http.NewRequest("POST", OLLAMA_QUERY_URI, bytes.NewReader(data))
	fmt.Println("sending request to ollama using URL", req.URL)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	resp, err := HttpClient().Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var event OllamaResponse

		if err := decoder.Decode(&event); err != nil {
			fmt.Println("Error decoding JSON response from ollama:", err)
			return "", err
		}
		ch <- event
		if event.Done {
			fmt.Println("done!")
			return "", nil
		}
	}
	return "", nil
}

func PingOllama() error {
	req, err := http.NewRequest("GET", OLLAMA_BASE_URI, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}
	resp, err := HttpClient().Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil

}
