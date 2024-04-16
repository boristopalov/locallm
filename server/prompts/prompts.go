package prompts

import (
	"fmt"

	"github.com/boristopalov/localsearch/utils"
)

func OllamaTemplate(q string) string {
	return fmt.Sprintf(`Question: %s`, q)
}

func Tools() map[string]utils.Tool {
	queryVectorDBTool := utils.Tool{
		Action: "QueryVectorDB",
		Description: `Useful for searching through added files and websites. 
		Search for keywords in the text, not whole questions. 
		Avoid relative words like "yesterday" and think about what could be in the text. 
		The input to this tool will be run against a vector db. The top results will be returned to you as json.`,
	}
	webSearchTool := utils.Tool{
		Action: "WebSearch",
		Description: `Useful for searching the internet. 
				You must use this tool if you're not 100% certain of the answer. 
				The top 10 results will be added to the vector db. 
				The top 3 results are also getting returned to you directly. 
				For more search queries through the same websites, use the VectorDB tool`,
	}
	// TODO: gotta be a better way of doing this
	tools := map[string]utils.Tool{
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

var SystemMessage = fmt.Sprintf(`
Assistant is a large language model.

Assistant is designed to be able to assist humans with a wide range of tasks, from answering simple questions to providing in-depth explanations and discussions on a wide range of topics. As a language model, Assistant is able to generate human-like text based on the input it receives, allowing it to engage in natural-sounding conversations and provide responses that are coherent and relevant to the topic at hand.

Assistant is constantly learning and improving, and its capabilities are constantly evolving. It is able to process and understand large amounts of text, and can use this knowledge to provide accurate and informative responses to a wide range of questions. Additionally, Assistant is able to generate its own text based on the input it receives, allowing it to engage in discussions and provide explanations and descriptions on a wide range of topics.

Overall, Assistant is a powerful tool that can help with a wide range of tasks and provide valuable insights and information on a wide range of topics. Whether you need help with a specific question or just want to have a conversation about a particular topic, Assistant is here to assist.

TOOLS:
------

Assistant has access to the following tools:

%s

To use a tool, please use the following format:

%s

You have to use your tools to answer questions. You may use tools more than once.

You MUST provide any sources / links you've used to answer the question.

Create your reply in the same language as the search string.

When you have a response to say to the Human, or if you do not need to use a tool, you MUST use the format:

%s

You MUST format your final answer (after Final Answer:) in markdown. 

Begin!
`,

	ToolDescriptions(), // list of tools
	fmt.Sprintf(
		`
	Thought: Do I need to use a tool? Yes
	Action: the action to take, should be one of [%s]
	Action Input: the input to the action
	Observation: the result of the action
	`, ToolActions()),
	`
	Thought: Do I need to use a tool? No
	Final Answer: [your response here]
	`,
)

func RephrasePrompt(history string, question string) string {
	var historyMessage = fmt.Sprintf(`Given the following conversation and a follow up question, rephrase the
	follow up question to be a standalone question.
	Chat History:
	%s
	Follow Up Input: %s
	Standalone question:`, history, question)
	return historyMessage
}
