package utils

import (
	"regexp"
	"strings"
)

type SplitOptions struct {
	MinLength  int
	MaxLength  int
	Overlap    int
	Splitter   string
	Delimiters string
}

func splitChunk(currChunks []string, maxLength, overlap int) (string, string, string) {
	chunkString := strings.Join(currChunks, " ")
	subChunk := chunkString[:maxLength]
	restChunk := chunkString[maxLength:]

	blankPosition := strings.IndexByte(restChunk, ' ')
	if blankPosition == -1 {
		blankPosition = strings.IndexByte(restChunk, '\n')
	}

	if blankPosition != -1 {
		subChunk += restChunk[:blankPosition]
		restChunk = restChunk[blankPosition:]
	}

	overlapText := ""

	if overlap > 0 {
		blankPosition = strings.LastIndexByte(subChunk[:len(subChunk)-overlap], ' ')
		if blankPosition == -1 {
			blankPosition = strings.LastIndexByte(subChunk[:len(subChunk)-overlap], '\n')
		}

		if blankPosition != -1 {
			overlapText = subChunk[blankPosition:]
		}
	}

	return subChunk, restChunk, overlapText
}

func Chunk(text string, options SplitOptions) []string {
	if options.MaxLength == 0 {
		options.MaxLength = 1000
	}

	if options.Splitter == "" {
		options.Splitter = "paragraph"
	}

	if options.Delimiters == "" {
		if options.Splitter == "sentence" {
			options.Delimiters = "([.!?\\n])\\s*"
		} else {
			options.Delimiters = "\\n{2,}"
		}
	}

	regex := options.Delimiters
	baseChunk := regexp.MustCompile(regex).Split(text, -1)

	var chunks []string
	var currChunks []string
	currChunkLength := 0

	for i := 0; i < len(baseChunk); i += 2 {
		subChunk := baseChunk[i]
		if i+1 < len(baseChunk) && baseChunk[i+1] != "" {
			subChunk += baseChunk[i+1]
		}

		currChunks = append(currChunks, subChunk)
		currChunkLength += len(subChunk)

		if currChunkLength >= options.MinLength {
			subChunk, restChunk, overlapText := splitChunk(currChunks, options.MaxLength, options.Overlap)
			chunks = append(chunks, subChunk)

			currChunks = nil
			currChunkLength = len(overlapText) + len(restChunk)

			if overlapText != "" {
				currChunks = append(currChunks, overlapText)
			}
			if restChunk != "" {
				currChunks = append(currChunks, restChunk)
			}
		}
	}

	if len(currChunks) > 0 {
		subChunk, restChunk, _ := splitChunk(currChunks, options.MaxLength, options.Overlap)
		if subChunk != "" {
			chunks = append(chunks, subChunk)
		}
		if restChunk != "" {
			chunks = append(chunks, restChunk)
		}
	}

	return chunks
}
