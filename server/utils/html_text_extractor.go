package utils

import (
	"strings"

	"golang.org/x/net/html"
)

func ExtractText(htmlContent string) string {
	var textBuilder strings.Builder
	tokenizer := html.NewTokenizer(strings.NewReader(htmlContent))

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		if tokenType == html.TextToken {
			text := strings.TrimSpace(tokenizer.Token().Data)
			if text != "" {
				textBuilder.WriteString(text)
				textBuilder.WriteString("\n")
			}
		}
	}

	return textBuilder.String()
}
