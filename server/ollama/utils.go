package ollama

import (
	"fmt"
	"net/http"
	"os"

	"github.com/boristopalov/localsearch/types"
	"github.com/boristopalov/localsearch/vars"
)

func PingOllama() error {
	resp, err := http.Get(vars.OLLAMA_BASE_URI)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func Tools() map[string]types.Tool {
	queryVectorDBTool := types.Tool{
		Action: "QueryVectorDB",
		Description: `Useful for running a similarity search on previously crawled websites. 
		Search for keywords in the text, not whole questions. 
		You MUST avoid relative words like "yesterday" in your action input as they will not help your input.
		The top results will be returned to you in JSON.`,
	}
	webSearchTool := types.Tool{
		Action: "WebSearch",
		Description: `Use this tool to search for websites that may answer your search query.
		The best websites (according to the search engine) are broken down into small parts and added to your vector database.
		The summaries of the top results will be returned to you in JSON.
		You can query the vector database later with other inputs to get other parts of these websites.`,
	}
	// TODO: gotta be a better way of doing this
	tools := map[string]types.Tool{
		"WebSearch":     webSearchTool,
		"QueryVectorDB": queryVectorDBTool,
	}
	return tools
}

func ToolDescriptions() string {
	tools := Tools()
	var toolDescriptions string
	for _, t := range tools {
		toolDescriptions += fmt.Sprintf("%s: %s,\n", t.Action, t.Description)
	}
	return toolDescriptions
}

func ToolActions() string {
	tools := Tools()
	var toolActions string
	for _, t := range tools {
		toolActions += t.Action + ","
	}
	return toolActions
}

func WriteToHistoryFile(text string) {
	fileName := "history/history.txt"
	// fileName := "history/history" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt" // e.x. history17591230.txt
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	} else {
		_, err := f.WriteString(text)
		if err != nil {
			fmt.Println(err)
		}
	}
	f.Close()
}
