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
		{"../../../file.html", "https://example.com/a/b/page.html", "https://example.com/file.html"},
		{"/file.html", "https://example.com/path/to/page.html", "https://example.com/file.html"},
		{"../file.html", "https://example.com/page.html", "https://example.com/file.html"},
		{"file.html", "https://example.com/path/to/", "https://example.com/path/to/file.html"},
		{"", "https://example.com/path/to/page.html", "https://example.com/path/to/page.html"},
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
		{"https://www.malgow.net", "www.malgow.net"},
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

func TestValidateLinks(t *testing.T) {
	type TC struct {
		name          string
		links         []string
		curURL        string
		domain        string
		allowExternal bool
		expected      []string
	}

	tests := []TC{
		{
			name:   "Relative links resolved correctly",
			curURL: "https://books.toscrape.com/",
			domain: "books.toscrape.com",
			links: []string{
				"index.html",
				"catalogue/category/books_1/index.html",
				"catalogue/category/books/travel_2/index.html",
			},
			allowExternal: false,
			expected: []string{
				"https://books.toscrape.com/",
				"https://books.toscrape.com/catalogue/category/books_1/",
				"https://books.toscrape.com/catalogue/category/books/travel_2/",
			},
		},
		{
			name:   "Reject external domains when allowExternal=false",
			curURL: "https://books.toscrape.com/",
			domain: "books.toscrape.com",
			links: []string{
				"http://google.com",
				"https://example.com",
				"catalogue/category/books/philosophy_7/",
			},
			allowExternal: false,
			expected: []string{
				"https://books.toscrape.com/catalogue/category/books/philosophy_7/",
			},
		},
		{
			name:   "Allow external domains when allowExternal=true",
			curURL: "https://books.toscrape.com/",
			domain: "books.toscrape.com",
			links: []string{
				"https://example.com",
				"catalogue/category/books/classics_6/index.html",
			},
			allowExternal: true,
			expected: []string{
				"https://example.com/",
				"https://books.toscrape.com/catalogue/category/books/classics_6/",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateLinks(tc.links, tc.curURL, tc.domain, tc.allowExternal)

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("ValidateLinks() failed\nExpected: %v\nGot: %v", tc.expected, got)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	type NormalizeTestCase struct {
		input    string
		expected string
	}
	tests := []NormalizeTestCase{
		{"https://example.com", "https://example.com/"},
		{"https://example.com/", "https://example.com/"},
		{"https://example.com/page#section", "https://example.com/page"},
		{"https://example.com/index.html", "https://example.com/"},
		{"https://example.com/foo/index.html", "https://example.com/foo/"},
		{"https://example.com//foo//bar", "https://example.com/foo/bar"},
		{"https://example.com//foo/index.html#top", "https://example.com/foo/"},
		{"https://example.com/foo/bar", "https://example.com/foo/bar"},
		{"ht!tp://bad-url", "ht!tp://bad-url"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i+1), func(t *testing.T) {
			got := NormalizeURL(test.input)
			if got != test.expected {
				t.Errorf("NormalizeURL(%q) = %q; want %q", test.input, got, test.expected)
			}
		})
	}
}
