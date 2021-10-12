package service

import (
	"context"
	"fmt"

	"github.com/twizar/teams/internal/domain/entity"
)

type teamRepository interface {
	All(ctx context.Context) ([]*entity.Team, error)
	Filter(ctx context.Context, minRating float64, leagues []string, orderBy string, limit int) ([]*entity.Team, error)
	ByIDs(ctx context.Context, ids []string) ([]*entity.Team, error)
}

type Teams struct {
	teamRepo teamRepository
}

func NewTeams(teamRepo teamRepository) *Teams {
	return &Teams{teamRepo: teamRepo}
}

func (t Teams) AllTeams(ctx context.Context) ([]*entity.Team, error) {
	teams, err := t.teamRepo.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all teams from repo error: %w", err)
	}

	return teams, nil
}

func (t Teams) FilterTeams(ctx context.Context, minRating float64, leagues []string, orderBy string, limit int) ([]*entity.Team, error) {
	teams, err := t.teamRepo.Filter(ctx, minRating, leagues, orderBy, limit)
	if err != nil {
		return nil, fmt.Errorf("getting all teams from repo error: %w", err)
	}

	return teams, nil
}

func (t Teams) TeamsByID(ctx context.Context, ids []string) ([]*entity.Team, error) {
	teams, err := t.teamRepo.ByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("getting teams by ID from repo error: %w", err)
	}

	return teams, nil
}
