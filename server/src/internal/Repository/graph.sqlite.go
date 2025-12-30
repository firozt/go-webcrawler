package repository

import (
	"database/sql"
	"log"
)

// represents node table
type PageNode struct {
	ID  uint64 `json:"id"`
	URL string `json:"url"`
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

func (repo GraphRepository) InsertPageNode(url string) {
	_, err := repo.db.Exec(`
		INSERT into pageNode(url) 
		VALUES (?)
	`, url)

	if err != nil {
		log.Fatal(err)
	}
}

func (repo GraphRepository) GetIDFromURL(url string) uint64 {
	var res uint64

	err := repo.db.QueryRow(`
        SELECT id FROM pageNode
        WHERE url = ?
    `, url).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("from URL not found: %s", url)
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
		log.Fatal(err)
	}
}

func (repo GraphRepository) GetAllLinkRelations(domain string) *LinkGraph {
	domainPattern := "%" + domain + "%"
	rows, err := repo.db.Query(`
		SELECT l.id, l.url, r.id, r.url
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
		var fromURL, toURL string
		if err := rows.Scan(&fromID, &fromURL, &toID, &toURL); err != nil {
			log.Fatal(err)
		}

		graph.Edges = append(graph.Edges, PageEdge{FromID: fromID, ToID: toID})

		if _, exists := seen[fromURL]; !exists {
			nodes = append(nodes, PageNode{ID: fromID, URL: fromURL})
			seen[fromURL] = true
		}
		if _, exists := seen[toURL]; !exists {
			nodes = append(nodes, PageNode{ID: toID, URL: toURL})
			seen[toURL] = true
		}
	}

	graph.Nodes = nodes
	return graph
}
