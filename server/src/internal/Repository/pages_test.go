package repository

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInsertPageFTS5(t *testing.T) {
	db := InitTestDB(t)
	defer db.Close()

	repo := NewPagesRepository(db)
	page := Page{
		URL:     "http://example.com",
		Title:   "Example Page",
		Content: "This is some test content",
	}

	err := repo.InsertPage(page)
	if err != nil {
		t.Fatalf("InsertPage failed: %v", err)
	}

	var url, title, content, crawledAt string
	err = repo.db.QueryRow(`SELECT url, title, content, crawled_at FROM pages WHERE url = ?`, page.URL).
		Scan(&url, &title, &content, &crawledAt)
	if err != nil {
		t.Fatalf("failed to query inserted page: %v", err)
	}

	if url != page.URL || title != page.Title || content != page.Content {
		t.Fatalf("row values do not match inserted page. got: %v %v %v", url, title, content)
	}
}

func TestSearchPages(t *testing.T) {
	db := InitTestDB(t)

	repo := NewPagesRepository(db)

	// Insert example page
	page := Page{
		URL:     "http://example.com",
		Title:   "Example Page",
		Content: "This is some test content",
	}

	if err := repo.InsertPage(page); err != nil {
		t.Fatalf("InsertPage failed: %v", err)
	}

	// Test case struct
	type SearchTest struct {
		query    string
		expected int
	}

	tests := []SearchTest{
		{"some", 1},
		{"this should return empty", 0},
		{"this is", 1},
	}

	for _, tc := range tests {
		t.Run("Test_keyword_'"+tc.query+"'", func(t *testing.T) {
			res := repo.SearchPages(tc.query, "", 10)

			if len(res) != tc.expected {
				t.Errorf("query %q: expected %d results, got %d",
					tc.query, tc.expected, len(res))
			}
		})
	}
	db.Close()
}
