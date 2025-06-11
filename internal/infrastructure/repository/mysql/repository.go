package mysql

import (
	"database/sql"
	"strings"

	"CardGameDB/internal/domain/card"
)

// Repository implements card.Repository using MySQL

type Repository struct {
	db *sql.DB
}

// New returns a new MySQL Repository
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Search implements card.Repository
func (r *Repository) Search(filter card.Filter) ([]card.Card, error) {
	query := "SELECT id, cost, upgrade_cost, faction, category, subcategory FROM cards"
	var where []string
	var args []interface{}

	if filter.ID != nil {
		where = append(where, "id = ?")
		args = append(args, *filter.ID)
	}
	if filter.Cost != nil {
		where = append(where, "cost = ?")
		args = append(args, *filter.Cost)
	}
	if filter.UpgradeCost != nil {
		where = append(where, "upgrade_cost = ?")
		args = append(args, *filter.UpgradeCost)
	}
	if filter.Faction != nil {
		where = append(where, "faction = ?")
		args = append(args, *filter.Faction)
	}
	if filter.Category != nil {
		where = append(where, "category = ?")
		args = append(args, *filter.Category)
	}
	if filter.SubCategory != nil {
		where = append(where, "subcategory = ?")
		args = append(args, *filter.SubCategory)
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []card.Card
	for rows.Next() {
		var c card.Card
		if err := rows.Scan(&c.ID, &c.Cost, &c.UpgradeCost, &c.Faction, &c.Category, &c.SubCategory); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

// Create inserts a new card into the database
func (r *Repository) Create(c card.Card) error {
	_, err := r.db.Exec(
		"INSERT INTO cards (id, cost, upgrade_cost, faction, category, subcategory) VALUES (?, ?, ?, ?, ?, ?)",
		c.ID, c.Cost, c.UpgradeCost, c.Faction, c.Category, c.SubCategory,
	)
	return err
}

// Update updates an existing card
func (r *Repository) Update(c card.Card) error {
	_, err := r.db.Exec(
		"UPDATE cards SET cost=?, upgrade_cost=?, faction=?, category=?, subcategory=? WHERE id=?",
		c.Cost, c.UpgradeCost, c.Faction, c.Category, c.SubCategory, c.ID,
	)
	return err
}
