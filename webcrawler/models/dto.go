package models

// BookWithAuthors is a data transfer object that contains both a book and its authors
type BookWithAuthors struct {
	Book    *Book
	Authors []string // Simple strings for author names
}
