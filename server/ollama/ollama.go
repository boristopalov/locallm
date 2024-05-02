package ollama

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/boristopalov/localsearch/tools"
	"github.com/boristopalov/localsearch/types"
	"github.com/boristopalov/localsearch/vars"
)

type ModelTextOutput struct {
	Text string
}

var memory []int
var history = ""
var currentQ = ""

/*
 */
func PromptModel(q string) (string, error) {
	query := types.OllamaQuery{
		Prompt:    q,
		KeepAlive: "5m",
		Model:     "llama3-custom",
		Context:   memory,
		System:    SystemMessage,
	}
	data, err := json.Marshal(query)
	if err != nil {
		fmt.Println("error marshalling data:", query)
		return "", err
	}
	fmt.Println(q)
	resp, err := http.Post(vars.OLLAMA_QUERY_URI, "application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var response strings.Builder
	for decoder.More() {
		var event types.OllamaResponse
		if err := decoder.Decode(&event); err != nil {
			fmt.Println("Error decoding JSON response from ollama:", err)
			return "", err
		}
		fmt.Print(event.Response)
		// fmt.Print(".")
		response.WriteString(event.Response)
		if event.Done {
			fmt.Println("done!")
			return response.String(), nil
		}
	}
	return "", nil
}

func getRephrasedPrompt(history string, question string) (string, error) {
	promptToRephrase := RephrasePromptTemplate(history, question)
	response, err := PromptModel(promptToRephrase)
	if err != nil {
		return "", errors.New("unable to rephrase prompt")
	}
	rephrasedPrompt := vars.RephrasedQuestionRegex.FindStringSubmatch(response)
	if len(rephrasedPrompt) > 1 {
		return rephrasedPrompt[1], nil
	}
	return "", fmt.Errorf("cannot find rephrased prompt regex match; prompt is: %s", response)
}

// helper function when determining the next action to take
func GetActionAndInput(s string) (string, string) {
	action := vars.ActionRegex.FindStringSubmatch(s)
	if len(action) > 1 {
		fmt.Println("found action to take...")
		actionInput := vars.ActionInputRegex.FindStringSubmatch(s)
		if len(actionInput) > 1 {
			fmt.Println("found action input:", actionInput[1])
			toolName := action[1]
			// make sure the model isn't hallucinating a tool
			_, ok := Tools()[toolName]
			if ok {
				fmt.Println("found tool with name", toolName)
				return toolName, actionInput[1]
			}
		}
		return "", ""
	}
	return "", ""
}

func AnswerQuestion(q string, clientDataChan chan types.DataStreamEvent) {
	response, err := PromptModel(q)
	if err != nil {
		fmt.Println("error in AnswerQuestion: ", err.Error())
		clientDataChan <- types.DataStreamEvent{
			EventName: "Error",
			Data:      "Model Error. Please try again.",
		}
		return
	}

	finalAnswer := vars.FinalAnswerRegex.FindStringSubmatch(response)
	fmt.Println("final answer regex result:", finalAnswer)
	if len(finalAnswer) > 1 {
		// fmt.Println("final answer: ", finalAnswer[1])
		clientDataChan <- types.DataStreamEvent{
			EventName: "Answer",
			Data:      finalAnswer[1],
		}
		ans := fmt.Sprintf("Answer: %s\n", finalAnswer[1])
		h := "Question: " + currentQ + "\n" + ans
		history += h
		WriteToHistoryFile(h)
		return
	}

	action, input := GetActionAndInput(response)
	fmt.Printf("got action {%s} and input {%s}\n", action, input)
	if action == "" || input == "" {
		clientDataChan <- types.DataStreamEvent{
			EventName: "Answer",
			Data:      response,
		}
		ans := fmt.Sprintf("Answer: %s\n", response)
		h := "Question: " + currentQ + "\n" + ans
		history += h
		WriteToHistoryFile(h)
		return
	}
	TakeAction(response, action, input, clientDataChan)
}

func TakeAction(prompt string, action string, input string, clientDataChan chan types.DataStreamEvent) {
	clientDataChan <- types.DataStreamEvent{
		EventName: action,
		Data:      input,
	}
	switch action {
	case "WebSearch":
		webSearchRes, err := tools.WebSearchTool.Execute(input)
		if err != nil {
			fmt.Println("error searching the web:", err.Error())
			clientDataChan <- types.DataStreamEvent{
				EventName: "Error",
				Data:      "Error searching the web. Please try again",
			}
			return
		}
		json, _ := json.Marshal(webSearchRes)
		clientDataChan <- types.DataStreamEvent{
			EventName: "WebSearchResult",
			Data:      string(json),
		}
		prompt = HumanTemplate(currentQ) + "\n" + prompt + fmt.Sprintf("Observation: %s\n", string(json))
		AnswerQuestion(prompt, clientDataChan)

	case "QueryVectorDB":
		queryDbRes, err := tools.QueryVectorDBTool.Execute(input)
		if err != nil {
			fmt.Println("error querying the DB", err.Error())
			clientDataChan <- types.DataStreamEvent{
				EventName: "Error",
				Data:      "Error querying the Vector DB. Please try again",
			}
			return
		}
		clientDataChan <- types.DataStreamEvent{
			EventName: "QueryVectorDBResult",
			Data:      queryDbRes,
		}
		prompt = HumanTemplate(currentQ) + "\n" + prompt + fmt.Sprintf("Observation: %s\n", queryDbRes)
		AnswerQuestion(prompt, clientDataChan)
	}

}

func StartAgentChain(question string, clientDataChan chan types.DataStreamEvent) {
	if len(history) > 1 {
		fmt.Println("rephrasing prompt...")
		newQuestion, err := getRephrasedPrompt(history, question)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			question = newQuestion
		}
	}
	currentQ = question
	AnswerQuestion(HumanTemplate(question), clientDataChan)
}
