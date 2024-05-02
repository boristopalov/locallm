package utils

import (
	"log"
	"strings"
	"unicode/utf8"

	"github.com/boristopalov/localsearch/vars"
)

// taken from https://github.com/tmc/langchaingo/
type RecursiveTextSplitter struct {
	Separators   []string
	ChunkSize    int
	MinChunkSize int
	ChunkOverlap int
	LenFunc      func(string) int
}

func DefaultRecursiveTextSplitter() RecursiveTextSplitter {
	return RecursiveTextSplitter{
		Separators:   []string{"\n\n", "\n", " ", ""},
		LenFunc:      utf8.RuneCountInString,
		ChunkSize:    vars.CHUNK_SIZE,
		MinChunkSize: vars.MIN_CHUNK_SIZE,
		ChunkOverlap: vars.CHUNK_OVERLAP,
	}
}

// SplitText splits a text into multiple text.
func (s RecursiveTextSplitter) SplitText(text string) ([]string, error) {
	finalChunks := make([]string, 0)

	// Find the appropriate separator
	separator := s.Separators[len(s.Separators)-1]
	newSeparators := []string{}
	for i, c := range s.Separators {
		if c == "" || strings.Contains(text, c) {
			separator = c
			newSeparators = s.Separators[i+1:]
			break
		}
	}

	splits := strings.Split(text, separator)
	goodSplits := make([]string, 0)

	// Merge the splits, recursively splitting larger texts.
	for _, split := range splits {
		if s.LenFunc(split) < s.ChunkSize {
			goodSplits = append(goodSplits, split)
			continue
		}

		if len(goodSplits) > 0 {
			mergedText := mergeSplits(goodSplits, separator, s.ChunkSize, s.ChunkOverlap, s.LenFunc)

			for _, chunk := range mergedText {
				// fmt.Println(chunk)
				if s.LenFunc(chunk) >= s.MinChunkSize {
					finalChunks = append(finalChunks, chunk)
				}
			}
			goodSplits = make([]string, 0)
		}

		if len(newSeparators) == 0 {
			finalChunks = append(finalChunks, split)
		} else {
			otherInfo, err := s.SplitText(split)
			if err != nil {
				return nil, err
			}
			finalChunks = append(finalChunks, otherInfo...)
		}
	}

	if len(goodSplits) > 0 {
		mergedText := mergeSplits(goodSplits, separator, s.ChunkSize, s.ChunkOverlap, s.LenFunc)
		finalChunks = append(finalChunks, mergedText...)
	}

	return finalChunks, nil
}

// joinDocs comines two documents with the separator used to split them.
func joinDocs(docs []string, separator string) string {
	return strings.TrimSpace(strings.Join(docs, separator))
}

func mergeSplits(splits []string, separator string, chunkSize int, chunkOverlap int, lenFunc func(string) int) []string { //nolint:cyclop
	docs := make([]string, 0)
	currentDoc := make([]string, 0)
	total := 0

	for _, split := range splits {
		totalWithSplit := total + lenFunc(split)
		if len(currentDoc) != 0 {
			totalWithSplit += lenFunc(separator)
		}

		maybePrintWarning(total, chunkSize)
		if totalWithSplit > chunkSize && len(currentDoc) > 0 {
			doc := joinDocs(currentDoc, separator)
			if doc != "" {
				docs = append(docs, doc)
			}

			for shouldPop(chunkOverlap, chunkSize, total, lenFunc(split), lenFunc(separator), len(currentDoc)) {
				total -= lenFunc(currentDoc[0]) //nolint:gosec
				if len(currentDoc) > 1 {
					total -= lenFunc(separator)
				}
				currentDoc = currentDoc[1:] //nolint:gosec
			}
		}

		currentDoc = append(currentDoc, split)
		total += lenFunc(split)
		if len(currentDoc) > 1 {
			total += lenFunc(separator)
		}
	}

	doc := joinDocs(currentDoc, separator)
	if doc != "" {
		docs = append(docs, doc)
	}

	return docs
}

func maybePrintWarning(total, chunkSize int) {
	if total > chunkSize {
		log.Printf(
			"[WARN] created a chunk with size of %v, which is longer then the specified %v\n",
			total,
			chunkSize,
		)
	}
}

// Keep popping if:
//   - the chunk is larger than the chunk overlap
//   - or if there are any chunks and the length is long
func shouldPop(chunkOverlap, chunkSize, total, splitLen, separatorLen, currentDocLen int) bool {
	docsNeededToAddSep := 2
	if currentDocLen < docsNeededToAddSep {
		separatorLen = 0
	}

	return currentDocLen > 0 && (total > chunkOverlap || (total+splitLen+separatorLen > chunkSize && total > 0))
}
