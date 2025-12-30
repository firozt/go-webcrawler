package repository

import (
	"testing"
)

func TestGetAllLinkRelations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// repo := NewGraphRepository(db)

}
