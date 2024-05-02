#!/bin/sh
docker compose down
# osascript -e 'tell app "Ollama" to quit' # gotta be a better why 

docker compose up -d


# ollama serve & 

echo "ollama running"
echo "spinning up main go server"
: > history.txt # clear chat history
go build && go install github.com/boristopalov/localsearch && ./localsearch


