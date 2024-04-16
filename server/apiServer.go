package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/boristopalov/localsearch/prompts"
	"github.com/boristopalov/localsearch/tools"
	"github.com/boristopalov/localsearch/utils"
)

var memory []int
var history = ""

func llamaPing(w http.ResponseWriter, r *http.Request) {
	err := utils.PingOllama()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Ollama server is reachable"))
}

func promptHandler(w http.ResponseWriter, r *http.Request) {
	// make sure SSE is working/set up
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE is required but not supported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")

	// decode incoming client request
	var userQuery utils.ClientQuery
	fmt.Println("decoding user query...")
	err := json.NewDecoder(r.Body).Decode(&userQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("sucessfully decoded user query")

	// channel to keep track of the state of the query
	// i.e. we need to tell the app when we have a final answer
	stateChan := make(chan string)

	if len(history) > 1 {
		fmt.Println("rephrasing prompt...")
		prompt, err := getRephrasedPrompt(history, userQuery)
		if err != nil {
			fmt.Println("cannot rephrase prompt")
			return
		}
		userQuery.Prompt = prompt
	}

	go startAgentChain(userQuery)

	// TODO: make this better
	// where should channels be closed?
	for event := range stateChan {
		if event == "DONE!" {
			// close(modelOutputChan)
			// close(stateChan)
			fmt.Fprint(w, "DONE!")
			flusher.Flush()
		}
	}
}

func getRephrasedPrompt(history string, userQuery utils.ClientQuery) (string, error) {
	promptToRephrase := prompts.RephrasePrompt(history, userQuery.Prompt)
	llmQuery := utils.OllamaQuery{
		Prompt:    promptToRephrase,
		KeepAlive: userQuery.KeepAlive,
		Model:     userQuery.Model,
		Context:   memory,
		System:    prompts.SystemMessage,
	}
	// the response from the model here
	// is just the rephrased prompt, NOT the answer to the question
	ch := make(chan utils.OllamaResponse)
	var prompt strings.Builder
	utils.PromptModel(llmQuery, ch)
	for event := range ch {
		// print out the incoming stream of text
		fmt.Print(event.Response)
		if event.Done {
			close(ch)
			fmt.Println("done reading response")
			// TODO: is this needed? i dont think so
			// memory = append(memory, event.Context...)
		}
		// keep track of the response
		prompt.WriteString(event.Response)
	}
	rephrasedPrompt := utils.RephrasedQuestionRegex.FindStringSubmatch(prompt.String())
	if len(rephrasedPrompt) > 1 {
		return rephrasedPrompt[1], nil
	}
	return "", fmt.Errorf("cannot find rephrased prompt regex matchm prompt is: %s", prompt.String())
}

func startModelOutputHandler(ch chan utils.OllamaResponse) {
	var responseBuilder strings.Builder

	for event := range ch {
		// print out the incoming stream of text
		fmt.Print(event.Response)
		if event.Done {
			close(ch)
			fmt.Println("done reading response")
			// TODO: is this needed? i dont think so
			// memory = append(memory, event.Context...)
		}
		// keep track of the response
		responseBuilder.WriteString(event.Response)
	}

	fmt.Println("response:", responseBuilder.String(), "-----------")

	// at this point responseBuilder has the full response of the model
	// check if the model either has a final answer,
	// or if it has determined it needs to take an action
	// if it needs to take an action,
	// use the tool it has determined it should use, and the tool input it has come up with
	finalAnswer := utils.FinalAnswerRegex.FindStringSubmatch(responseBuilder.String())
	if len(finalAnswer) > 1 {
		fmt.Println("ANSWER!", finalAnswer[1])
		return
	}

	action := utils.ActionRegex.FindStringSubmatch(responseBuilder.String())
	fmt.Println("action array", action)
	fmt.Println("action regex find:", action[1])
	if len(action) > 1 {
		fmt.Println("found action to take...")
		actionInput := utils.ActionInputRegex.FindStringSubmatch(responseBuilder.String())
		if len(actionInput) > 1 {
			fmt.Println("found action input:", actionInput[1])
			toolName := action[1]
			// make sure the model isn't hallucinating a tool
			_, ok := prompts.Tools()[toolName]
			if ok {
				fmt.Println("found tool with name", toolName)
				switch toolName {
				case "WebSearch":
					webSearchRes, err := tools.WebSearchTool_.Execute(actionInput[1])
					if err != nil {
						fmt.Println("error searching the web:", err.Error())
						return
					}
					responseBuilder.WriteString(fmt.Sprintf("Observation: %s\n", webSearchRes))
					// history += fmt.Sprintf("Q: %s\nA: %s", rephrasedPrompt[1], webSearchRes)

				case "QueryVectorDB":
					vectoryQuery, err := utils.GetTextEmbedding(actionInput[1])
					if err != nil {
						fmt.Println("error getting text embedding:", err.Error())
						return
					}
					queryDbRes, err := tools.QueryVectorDBTool_.Execute(vectoryQuery)
					if err != nil {
						fmt.Println("error querying the DB", err.Error())
						return
					}
					responseBuilder.WriteString(fmt.Sprintf("Observation: %s\n", queryDbRes))
					// history += fmt.Sprintf("Q: %s\nA: %s", rephrasedPrompt[1], queryDbRes)
				}
				cq := utils.ClientQuery{
					Prompt:    responseBuilder.String(),
					Model:     "llama2",
					KeepAlive: "5m",
				}
				// query the model with the observation
				//TODO: fix function naming
				startAgentChain(cq)
				return
			}
		}
	}
}

// TODO: front-end relies on the name of the event here. Should be something like text-data
func startAgentChain(q utils.ClientQuery) {
	// channel that will stream the output of the model
	ch := make(chan utils.OllamaResponse)
	go startModelOutputHandler(ch)
	llmQuery := utils.OllamaQuery{
		Prompt:    q.Prompt,
		KeepAlive: q.KeepAlive,
		Model:     q.Model,
		Context:   memory,
		System:    prompts.SystemMessage,
	}
	fmt.Println("incoming client request:", llmQuery.Prompt)
	_, err := utils.PromptModel(llmQuery, ch)
	if err != nil {
		fmt.Println(err)
	}
}

func startServer() {
	http.HandleFunc("/ping", llamaPing)
	http.HandleFunc("/prompt", promptHandler)
	http.ListenAndServe(":8000", nil)

}
