package converter

import (
	"github.com/twizar/common/pkg/dto"
	"github.com/twizar/teams/internal/domain/entity"
)

func DTOtoEntity(t dto.Team) *entity.Team {
	return entity.NewTeam(t.ID, t.Name, t.League, t.Rating)
}

func DTOsToEntities(dtoTeams []dto.Team) []*entity.Team {
	entityTeams := make([]*entity.Team, len(dtoTeams))
	for i, team := range dtoTeams {
		entityTeams[i] = DTOtoEntity(team)
	}

	return entityTeams
}

func EntitiesToDTOs(entityTeams []*entity.Team) []dto.Team {
	dtoTeams := make([]dto.Team, len(entityTeams))
	for i, team := range entityTeams {
		dtoTeams[i] = dto.Team{ID: team.ID(), Name: team.Name(), Rating: team.Rating(), League: team.League()}
	}

	return dtoTeams
}
