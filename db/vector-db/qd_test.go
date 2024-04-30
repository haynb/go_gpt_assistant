package drant_db

import (
	"testing"
)

func TestDeleteCollection(t *testing.T) {
	DeleteCollection("test")
}

func TestCreateCollection(t *testing.T) {
	CreateCollection("test")
}
