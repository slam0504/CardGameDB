package mysql

import (
	"CardGameDB/internal/domain/card"
	"gorm.io/gorm"
)

// Repository implements card.Repository using GORM with MySQL

type Repository struct {
	db *gorm.DB
}

// New returns a new Repository backed by GORM
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Search implements card.Repository
func (r *Repository) Search(filter card.Filter) ([]card.Card, error) {
	var cards []card.Card
	q := r.db.Model(&card.Card{})

	if filter.ID != nil {
		q = q.Where("id = ?", *filter.ID)
	}
	if filter.Cost != nil {
		q = q.Where("cost = ?", *filter.Cost)
	}
	if filter.UpgradeCost != nil {
		q = q.Where("upgrade_cost = ?", *filter.UpgradeCost)
	}
	if filter.Faction != nil {
		q = q.Where("faction = ?", *filter.Faction)
	}
	if filter.Category != nil {
		q = q.Where("category = ?", *filter.Category)
	}
	if filter.SubCategory != nil {
		q = q.Where("subcategory = ?", *filter.SubCategory)
	}

	err := q.Find(&cards).Error
	return cards, err
}

// Create inserts a new card into the database
func (r *Repository) Create(c card.Card) error {
	return r.db.Create(&c).Error
}

// Update updates an existing card
func (r *Repository) Update(c card.Card) error {
	return r.db.Model(&card.Card{}).Where("id = ?", c.ID).Updates(c).Error
}
