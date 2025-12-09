package webcrawler

import (
	"testing"
)

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
