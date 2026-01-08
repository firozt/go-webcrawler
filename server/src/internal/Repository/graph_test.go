package repository

import (
	"testing"
)

func TestGetAllLinkRelations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewGraphRepository(db)
	// url_1 -> url_2 -> url_3
	repo.InsertPageNode("https://test.com/1", "test_title_1")
	repo.InsertPageNode("https://test.com/2", "test_title_2")
	repo.InsertPageNode("https://test.com/3", "test_title_3")
	repo.InsertPageEdge("https://test.com/1", "https://test.com/2")
	repo.InsertPageEdge("https://test.com/2", "https://test.com/3")

	expected := LinkGraph{
		Nodes: []PageNode{
			{
				ID:    1,
				URL:   "https://test.com/1",
				Title: "test_title_1",
			},
			{
				ID:    2,
				URL:   "https://test.com/2",
				Title: "test_title_2",
			},
			{
				ID:    3,
				URL:   "https://test.com/3",
				Title: "test_title_3",
			},
		},
		Edges: []PageEdge{
			{
				FromID: 1,
				ToID:   2,
			},
			{
				FromID: 2,
				ToID:   3,
			},
		},
	}

	t.Run("Test", func(t *testing.T) {
		res := repo.GetAllLinkRelations("https://www.test.com")
		if len(res.Nodes) != len(expected.Nodes) || len(res.Edges) != len(expected.Edges) {
			t.Fatalf("Length of nodes or edges of actual do not match expected")
		}

		// get urls into a list for expected and actual
		actualUrls := make([]string, len(res.Nodes))
		for i, node := range res.Nodes {
			actualUrls[i] = node.URL
		}

		expectedUrls := make([]string, len(expected.Nodes))
		for i, node := range res.Nodes {
			expectedUrls[i] = node.URL
		}

		// compare urls
		for _, url := range actualUrls {
			if !containsString(expectedUrls, url) {
				t.Fatalf("Node list urls are not the same")
			}
		}
	})
}

func containsString(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
