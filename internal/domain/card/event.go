package card

// SearchRequested event triggered when a search is requested.
type SearchRequested struct {
	Filter Filter
	Reply  chan SearchResult
}

// SearchResult event contains search results.
type SearchResult struct {
	Cards []Card
	Err   error
}

// CreateRequested event triggered when a create is requested.
type CreateRequested struct {
	Card  Card
	Reply chan error
}

// UpdateRequested event triggered when an update is requested.
type UpdateRequested struct {
	Card  Card
	Reply chan error
}
