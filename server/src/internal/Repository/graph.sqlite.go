package repository

import (
	"database/sql"
	"fmt"
	"log"
)

// represents node table
type PageNode struct {
	ID    uint64 `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

// represent many to many relationship between pages
type PageEdge struct {
	FromID uint64 `json:"from_id"`
	ToID   uint64 `json:"to_id"`
}

// combines edges and nodes to create a graph object
type LinkGraph struct {
	Nodes []PageNode `json:"nodes"`
	Edges []PageEdge `json:"edges"`
}

// repository to control db
type GraphRepository struct {
	db *sql.DB
}

// constructor
func NewGraphRepository(db *sql.DB) *GraphRepository {
	return &GraphRepository{db: db}
}

func (repo GraphRepository) InsertPageNode(url string, title string) {
	// prevents duplicate url (candidate key)
	if repo.isRowExist(url) {
		fmt.Printf("--stopping insert (duplicate) for %s\n", url)
		return
	}
	_, err := repo.db.Exec(`
        INSERT INTO pageNode(url, title)
        VALUES (?, ?)
    `, url, title)

	if err != nil {
		log.Fatal(err)
	}
}

func (repo GraphRepository) isRowExist(url string) bool {
	var count uint16
	err := repo.db.QueryRow(`
		SELECT count(*) FROM pageNode
		WHERE url = ?
	`, url).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func (repo GraphRepository) GetIDFromURL(url string) uint64 {
	var res uint64

	err := repo.db.QueryRow(`
        SELECT id FROM pageNode
        WHERE url = ?
    `, url).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatal("No rows matched")
		}
		log.Fatal(err)
	}
	return res
}

func (repo GraphRepository) InsertPageEdge(fromURL string, toURL string) {
	// get ID from fromURL and toURL
	fromID := repo.GetIDFromURL(fromURL)
	toID := repo.GetIDFromURL(toURL)

	// insert the edge
	_, err := repo.db.Exec(`
        INSERT INTO pageLink(from_id, to_id)
        VALUES (?, ?)
    `, fromID, toID)
	if err != nil {
		log.Fatal("Unable to insert new page to graph database - ", err)
	}
}
func (repo GraphRepository) AlterPageNodeTitle(url string, title string) {
	id := repo.GetIDFromURL(url)

	_, err := repo.db.Exec(`
        UPDATE pageNode
        SET title = ?
        WHERE id = ?;
    `, title, id)

	if err != nil {

		log.Fatal("Attempting to alter a non existing page node - ", err)
	}
}

func (repo GraphRepository) GetAllLinkRelations(domain string) *LinkGraph {
	domainPattern := "%" + domain + "%"
	rows, err := repo.db.Query(`
		SELECT l.id, l.title , l.url, r.id, r.title , r.url
		FROM pageLink
		JOIN pageNode AS l ON l.id = pageLink.from_id
		JOIN pageNode AS r ON r.id = pageLink.to_id
		WHERE l.url LIKE ? AND r.url LIKE ?
	`, domainPattern, domainPattern)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	graph := &LinkGraph{}
	nodes := []PageNode{}
	seen := make(map[string]bool)

	for rows.Next() {
		var fromID, toID uint64
		var fromURL, toURL, fromTitle, toTitle string
		if err := rows.Scan(&fromID, &fromTitle, &fromURL, &toID, &toTitle, &toURL); err != nil {
			log.Fatal(err)
		}

		graph.Edges = append(graph.Edges, PageEdge{FromID: fromID, ToID: toID})

		if _, exists := seen[fromURL]; !exists {
			nodes = append(nodes, PageNode{ID: fromID, URL: fromURL, Title: fromTitle})
			seen[fromURL] = true
		}
		if _, exists := seen[toURL]; !exists {
			nodes = append(nodes, PageNode{ID: toID, URL: toURL, Title: toTitle})
			seen[toURL] = true
		}
	}

	graph.Nodes = nodes
	return graph
}
