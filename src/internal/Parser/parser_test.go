package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
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

type DFSTestCase struct {
	html            string
	expectedLinks   []string
	expectedContent []string
}

func TestDFS(t *testing.T) {
	// inputs
	testCases := []DFSTestCase{
		{
			`<div>Hello <a href="https://example.com">world</a></div>`,
			[]string{"https://example.com"},
			[]string{"Hello ", "world"},
		},
		{
			`<div>Start <div>Nested <a href="https://nested.com/page">link</a></div> End</div>`,
			[]string{"https://nested.com/page"},
			[]string{"Start ", "Nested ", "link", " End"},
		},
		{
			`<p><a href="https://one.com">One</a> and <a href="https://two.com">Two</a></p>`,
			[]string{"https://one.com", "https://two.com"},
			[]string{"One", " and ", "Two"},
		},
		{
			`<span>Just some text without links.</span>`,
			[]string{},
			[]string{"Just some text without links."},
		},
	}
	for i, testcase := range testCases {
		doc, err := html.Parse(strings.NewReader(testcase.html))

		if err != nil {
			t.Errorf("Unexpected error in test %v", err)
		}

		link, content := []string{}, []string{}
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			dfs(doc, &content, &link)

			if !reflect.DeepEqual(link, testcase.expectedLinks) {
				t.Errorf("wanted %v got %v", testcase.expectedLinks, link)
			}
			if !reflect.DeepEqual(content, testcase.expectedContent) {
				t.Errorf("wanted %v got %v", testcase.expectedContent, content)
			}
		})
	}
}
