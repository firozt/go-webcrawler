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

func TestAbsolutePathToUrl(t *testing.T) {
	testCases := [][3]string{
		// Root-relative paths (absolute within the same domain)
		{"/index.html", "https://example.com/path/to/page.html", "https://example.com/index.html"},
		{"/images/logo.png", "https://example.com/path/to/page.html", "https://example.com/images/logo.png"},
		{"/about/us.html", "https://example.com/path/to/page.html", "https://example.com/about/us.html"},
		{"/", "https://example.com/path/to/page.html", "https://example.com/"},

		// Root-relative edge case
		{"/docs/", "https://example.com/path/to/page.html", "https://example.com/docs/"},
	}
	for i, testCase := range testCases {
		absPath, curPath, expected := testCase[0], testCase[1], testCase[2]

		t.Run(fmt.Sprintf("case_%d_%s", i, absPath), func(t *testing.T) {
			actual, err := absolutePathToUrl(absPath, curPath)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if actual != expected {
				t.Errorf("got %v, want %v", actual, expected)
			}
		})
	}
}
