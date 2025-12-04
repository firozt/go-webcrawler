package parser

import (
	"fmt"
	"testing"
)

func TestRelativePathToUrl(t *testing.T) {
	testCases := [][3]string{
		{"file.html", "https://example.com/path/to/page.html", "https://example.com/path/to/file.html"},
		{"../file.html", "https://example.com/path/to/page.html", "https://example.com/path/file.html"},
		{"../../file.html", "https://example.com/a/b/c/page.html", "https://example.com/a/file.html"},
		{"./file.html", "https://example.com/path/to/page.html", "https://example.com/path/to/file.html"},
		{".", "https://example.com/path/to/page.html", "https://example.com/path/to/"},
		{"../sibling/file.html", "https://example.com/path/to/page.html", "https://example.com/path/sibling/file.html"},
	}

	for i, testCase := range testCases {
		relPath, curPath, expected := testCase[0], testCase[1], testCase[2]

		t.Run(fmt.Sprintf("case_%d_%s", i, relPath), func(t *testing.T) {
			actual, err := relativePathToUrl(relPath, curPath)
			if err != nil || actual != expected {
				t.Errorf("got %v, want %v", actual, expected)
			}
		})
	}
}
