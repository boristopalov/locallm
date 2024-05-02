package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boristopalov/localsearch/ollama"
	"github.com/boristopalov/localsearch/types"
	"github.com/boristopalov/localsearch/utils"
)

var clientDataChan = make(chan types.DataStreamEvent)

func llamaPing(w http.ResponseWriter, r *http.Request) {
	err := ollama.PingOllama()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Ollama server is reachable"))
}

func promptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	if r.Method != "POST" {
		w.Write([]byte("This endpoint only supports POST requests"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// decode incoming client request
	var userQuery types.ClientQuery // not much need for this rn
	err := json.NewDecoder(r.Body).Decode(&userQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go ollama.StartAgentChain(userQuery.Prompt, clientDataChan)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("sending ya prompt over"))
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	// make sure SSE is working/set up
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE is required but not supported", http.StatusInternalServerError)
		return
	}

	for event := range clientDataChan {
		_, err := fmt.Fprint(w, utils.FormatSSEMessage(event))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			flusher.Flush()
			fmt.Println("data has been sent over SSE")
		}
	}
}

func startServer() {
	http.HandleFunc("/ping", llamaPing)
	http.HandleFunc("/prompt", promptHandler)
	http.HandleFunc("/stream", streamHandler)
	http.ListenAndServe(":8000", nil)
}
