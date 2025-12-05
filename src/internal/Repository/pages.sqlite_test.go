package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:") // in-memory DB
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}

	// creates similar table as prod (cant do fts5)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			title TEXT,
			content TEXT,
			crawled_at TEXT
		);
    `)
	if err != nil {
		t.Fatalf("failed to create FTS5 pages table: %v", err)
	}

	return db
}

func TestInsertPageFTS5(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPagesRepository(db)
	// creates page obj
	page := Page{
		URL:     "http://example.com",
		Title:   "Example Page",
		Content: "This is some test content",
	}

	// insert the page into the db
	err := repo.InsertPage(page)
	if err != nil {
		t.Fatalf("InsertPage failed: %v", err)
	}

	// verify
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
