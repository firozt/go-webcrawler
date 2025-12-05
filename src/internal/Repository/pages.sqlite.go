/*
This file handles all the interactions with the database to do with the pages virtual
table.
*/
package repository

import (
	"database/sql"
	"log"
	"time"
)

type PagesRepository struct {
	db *sql.DB
}

// struct to represent a row in the pages table
type Page struct {
	Content string
	URL     string
	Title   string
}

// constructor for this object given a pre existing DB object
func NewPagesRepository(db *sql.DB) *PagesRepository {
	return &PagesRepository{db: db}
}

func (p PagesRepository) InsertPage(page Page) error {
	_, err := p.db.Exec(`
        INSERT INTO pages (url, title, content, crawled_at)
        VALUES (?, ?, ?, ?)
    `, page.URL, page.Title, page.Content, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

// searches for phrase from DB
func (p PagesRepository) SearchPages(phrase string) []Page {
	var res []Page
	rows, err := p.db.Query(`
		SELECT * FROM pages
		WHERE pages MATCH %s
		LIMIT 10
	`)
	if err != nil {
		return res
	}

	defer rows.Close()

	for rows.Next() {
		var url, title, content string
		err := rows.Scan(&url, &title, &content)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, Page{
			URL:     url,
			Title:   title,
			Content: content,
		})
	}

	return res
}
