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
		{"B-B0TheDeliveryBot", "https://malgow.net", "https://malgow.net/B-B0TheDeliveryBot"},

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
			[]string{"hello ", "world"},
		},
		{
			`<div>Start <div>Nested <a href="https://nested.com/page">link</a></div> End</div>`,
			[]string{"https://nested.com/page"},
			[]string{"start ", "nested ", "link", " end"},
		},
		{
			`<p><a href="https://one.com">One</a> and <a href="https://two.com">Two</a></p>`,
			[]string{"https://one.com", "https://two.com"},
			[]string{"one", " and ", "two"},
		},
		{
			`<span>Just some text without links.</span>`,
			[]string{},
			[]string{"just some text without links."},
		},
	}
	for i, testcase := range testCases {
		doc, err := html.Parse(strings.NewReader(testcase.html))

		if err != nil {
			t.Errorf("Unexpected error in test %v", err)
		}

		link, content := []string{}, []string{}
		var title string
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			dfs(doc, &content, &link, &title)

			if !reflect.DeepEqual(link, testcase.expectedLinks) {
				t.Errorf("wanted %v got %v", testcase.expectedLinks, link)
			}
			if !reflect.DeepEqual(content, testcase.expectedContent) {
				t.Errorf("wanted %v got %v", testcase.expectedContent, content)
			}
		})
	}
}

func TestGetDomain(t *testing.T) {
	type TestGetDomainInput struct {
		url      string
		expected string
	}

	testCases := []TestGetDomainInput{
		{"https://sub.example.com:8080/path?query=1", "sub.example.com"},
		{"http://example.org/abc", "example.org"},
		{"https://another.example.net", "another.example.net"},
		{"http://localhost:3000/test", "localhost"},
		{"https://example.com", "example.com"},
		{"ftp://example.com", ""},
		{"", ""},
	}

	for _, test := range testCases {
		t.Run(test.url, func(t *testing.T) {
			got := GetDomain(test.url)
			if got != test.expected {
				t.Errorf("getDomain(%q) = %q; want %q", test.url, got, test.expected)
			}
		})
	}
}
