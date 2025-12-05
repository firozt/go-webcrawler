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

	// creates FTS5 table
	_, err = db.Exec(`
        CREATE VIRTUAL TABLE IF NOT EXISTS pages USING fts5(
            url,
            title,
            content,
            crawled_at,
            tokenize='porter'
        )
    `)

	_, _ = db.Exec(`
		DELETE from pages;
	`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
