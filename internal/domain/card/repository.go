package card

// Repository defines card persistence
// Search by fields

type Repository interface {
	// Search returns cards that match the filters
	Search(filter Filter) ([]Card, error)

	// Create stores a new card
	Create(c Card) error

	// Update modifies an existing card
	Update(c Card) error
}

// Filter defines search parameters
// all fields are optional; zero or empty values mean not filtering

type Filter struct {
	ID          *int
	Cost        *int
	UpgradeCost *int
	Faction     *string
	Category    *string
	SubCategory *string
}
