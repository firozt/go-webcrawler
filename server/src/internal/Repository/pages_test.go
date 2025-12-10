package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// starts in memory DB for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE VIRTUAL TABLE pages USING fts5(
			url, title, content, crawled_at
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestInsertPageFTS5(t *testing.T) {
	db := setupTestDB(t)
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
	err = db.QueryRow(`SELECT url, title, content, crawled_at FROM pages WHERE url = ?`, page.URL).
		Scan(&url, &title, &content, &crawledAt)
	if err != nil {
		t.Fatalf("failed to query inserted page: %v", err)
	}

	if url != page.URL || title != page.Title || content != page.Content {
		t.Fatalf("row values do not match inserted page. got: %v %v %v", url, title, content)
	}
}

func TestSearchPages(t *testing.T) {
	db := InitDB()
	defer db.Close()

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
		t.Run("Test "+tc.query, func(t *testing.T) {
			t.Parallel()
			res := repo.SearchPages(tc.query, 10)

			if len(res) != tc.expected {
				t.Errorf("query %q: expected %d results, got %d",
					tc.query, tc.expected, len(res))
			}
		})
	}
}
