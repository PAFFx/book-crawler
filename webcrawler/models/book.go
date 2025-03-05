package models

type Book struct {
	ID          uint   `gorm:"primaryKey"`      // Primary key
	HTMLHash    string `gorm:"unique;not null"` // Unique hash of the HTML content
	URL         string `gorm:"not null"`        // URL of the product page may not be unique
	ImageURL    string // URL of the product image
	Title       string // Title of the product
	ISBN        string // ISBN of the product
	Description string // Description of the product
	IsIndexed   bool   `gorm:"default:false;not null"` // Whether the book has already been indexed
}

type Author struct {
	ID     uint   `gorm:"primaryKey"`
	BookID uint   // Foreign key
	Book   Book   `gorm:"foreignKey:BookID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Name   string `gorm:"not null"`
}
