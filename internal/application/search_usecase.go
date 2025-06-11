package application

import "CardGameDB/internal/domain/card"

// SearchUseCase handles card searching
func NewSearchUseCase(repo card.Repository) *SearchUseCase {
	return &SearchUseCase{repo: repo}
}

// SearchUseCase respond to SearchRequested events

type SearchUseCase struct {
	repo card.Repository
}

// Query performs a search and returns the cards
func (uc *SearchUseCase) Query(filter card.Filter) ([]card.Card, error) {
	return uc.repo.Search(filter)
}

// Handle handles SearchRequested events and sends SearchResult
func (uc *SearchUseCase) Handle(event card.SearchRequested) {
	cards, err := uc.repo.Search(event.Filter)
	event.Reply <- card.SearchResult{Cards: cards, Err: err}
}
