package application

import "CardGameDB/internal/domain/card"

// ManageUseCase handles creating and updating cards

type ManageUseCase struct {
	repo card.Repository
}

func NewManageUseCase(repo card.Repository) *ManageUseCase {
	return &ManageUseCase{repo: repo}
}

// HandleCreate handles CreateRequested events
func (uc *ManageUseCase) HandleCreate(event card.CreateRequested) {
	err := uc.repo.Create(event.Card)
	event.Reply <- err
}

// HandleUpdate handles UpdateRequested events
func (uc *ManageUseCase) HandleUpdate(event card.UpdateRequested) {
	err := uc.repo.Update(event.Card)
	event.Reply <- err
}
