package utils_test

import (
	"os"
	"testing"

	"github.com/boristopalov/localsearch/utils"
)

// Define a test case struct to hold input and expected output
type testCase struct {
	name     string
	html     string
	expected string
}

func TestExtractText(t *testing.T) {
	// Define test cases
	d, err := os.ReadFile("../html_test.txt")
	if err != nil {
		t.Error("error reading file: ", err.Error())
		return
	}

	testCases := []testCase{
		{
			name:     "Simple HTML content",
			html:     "<div><p>This is some <b>bold</b> and <i>italic</i> text.</p><p>Another paragraph.</p></div>",
			expected: "This is some bold and italic text. Another paragraph.",
		},
		{
			name:     "Empty HTML content",
			html:     "",
			expected: "",
		},
		{
			name:     "HTML content containing only whitespace",
			html:     "    ",
			expected: "",
		},
		{
			name:     "HTML content containing only comments",
			html:     "<!-- This is a comment -->",
			expected: "",
		},
		{
			name:     "HTML content containing special characters",
			html:     "<div>&lt; &gt; &amp; &quot; &apos;</div>",
			expected: "< > & \" '",
		},
		{
			name:     "long HTML with table",
			html:     string(d),
			expected: "",
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute the function under test
			result := utils.ExtractText(tc.html)

			// Compare the result with expected output
			if result != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, result)
			}
		})
	}
}
