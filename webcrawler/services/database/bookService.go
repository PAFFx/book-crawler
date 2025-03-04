package database

import (
	"book-search/webcrawler/models"
	"errors"

	"gorm.io/gorm"
)

// StoreBookWithAuthors stores a book and its authors in the database
func StoreBookWithAuthors(db *gorm.DB, bookWithAuthors *models.BookWithAuthors) error {
	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the book first
	if err := tx.Create(bookWithAuthors.Book).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create authors with references to the book
	for _, authorName := range bookWithAuthors.Authors {
		author := &models.Author{
			BookID: bookWithAuthors.Book.ID,
			Name:   authorName,
		}
		if err := tx.Create(author).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func StoreBook(db *gorm.DB, book *models.Book) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&book).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func CheckBookExists(db *gorm.DB, htmlHash string) (bool, error) {
	// Validate input
	if htmlHash == "" {
		return false, errors.New("html hash cannot be empty")
	}

	// Check if table exists first
	if !db.Migrator().HasTable(&models.Book{}) {
		return false, errors.New("books table does not exist")
	}

	// Use a more efficient count query that returns early
	var count int64
	err := db.Model(&models.Book{}).
		Where("html_hash = ?", htmlHash).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
