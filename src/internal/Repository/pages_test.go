package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

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

func TestSearchPageFTS5(t *testing.T) {
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

	// Test 1: match should return 1 result
	res := repo.SearchPages("some")
	if len(res) != 1 {
		t.Errorf("expected 1 result, got %d", len(res))
	} else {
		t.Logf("Test 1 passed: found page with title %q", page.Title)
	}

	// Test 2: no match should return 0 results
	res = repo.SearchPages("this should be empty")
	if len(res) != 0 {
		t.Errorf("expected 0 results, got %d", len(res))
	} else {
		t.Log("Test 2 passed: no results returned")
	}
}
