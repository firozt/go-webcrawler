/*
This file creates the DB (if not existant) and returns the DB object
back to the caller.

The files only purpose is to initialise the DB
*/
package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const PATH_TO_DB string = "data/mydb.sqlite"

func InitDB() *sql.DB {
	// opens DB file
	db, err := sql.Open("sqlite3", PATH_TO_DB)
	if err != nil {
		log.Fatal(err)
	}

	createPagesTable(db)
	createRelationshipsTable(db)

	return db
}

func createRelationshipsTable(db *sql.DB) {
	// creates node table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS pageNode(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
	_, _ = db.Exec(`
		DELETE from pageNode;
	`)

	// creates links table (many to many)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pageLink(
			from_id INTEGER,
			to_id INTEGER
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
	_, _ = db.Exec(`
		DELETE from pageLink;
	`)
}

func createPagesTable(db *sql.DB) {
	// creates FTS5 table
	_, err := db.Exec(`
        CREATE VIRTUAL TABLE IF NOT EXISTS pages USING fts5(
            url,
            title,
            content,
            crawled_at,
            tokenize='porter'
        )
    `)

	// removes pre existing (if exists)
	_, _ = db.Exec(`
		DELETE from pages;
	`)

	if err != nil {
		log.Fatal(err)
	}

}
