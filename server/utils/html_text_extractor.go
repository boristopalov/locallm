package utils

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func ExtractText(html string) string {
	policy := bluemonday.StrictPolicy()
	result := policy.Sanitize(html)
	// result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, "\t", "")
	result = strings.ReplaceAll(result, "&#34;", "")
	result = strings.ReplaceAll(result, "&#39;", "")
	return result
}
