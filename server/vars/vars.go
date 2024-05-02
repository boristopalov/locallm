package vars

import (
	"fmt"
	"regexp"
)

const OLLAMA_PORT = 11434

// var OLLAMA_PORT = os.Getenv("OLLAMA_HOST")
const LOCALHOST = "http://127.0.0.1"

var OLLAMA_BASE_URI = fmt.Sprintf("%s:%d", LOCALHOST, OLLAMA_PORT)
var OLLAMA_QUERY_URI = OLLAMA_BASE_URI + "/api/generate"
var OLLAMA_EMBED_URI = OLLAMA_BASE_URI + "/api/embeddings"

// var OLLAMA_START_TEMPLATE = ""
const SEARXNG_PORT = 8080

var SEARXNG_BASE_URI = fmt.Sprintf("%s:%d", LOCALHOST, SEARXNG_PORT)
var SEARXNG_SEARCH_URI = SEARXNG_BASE_URI + "/search"

var ActionRegex = regexp.MustCompile(`Action: (.*)`)
var ActionInputRegex = regexp.MustCompile(`Action Input: "?(.*)"?`)
var FinalAnswerRegex = regexp.MustCompile(`Final Answer:\s*(.*[\s\S]*)`)
var RephrasedQuestionRegex = regexp.MustCompile(`Standalone question:\s*(.*[\s\S]*)`)

const CHUNK_SIZE = 1000
const CHUNK_OVERLAP = 200
const MIN_CHUNK_SIZE = 400

var CHUNK_SEPARATORS = [4]string{"\n\n", "\n", " ", ""}

// var CONTEXT_SIZE = 8192
const EMBEDDINGS_MODEL = "nomic-embed-text"
const COLLECTION_NAME = "localsearch"
const DB_NAME = "localsearch"
const EMBEDDING_DIMS = 768

const MAX_ATTEMPTS = 5
