package public

import (
	context "silent/src/internal/context/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/dto"
	"silent/src/internal/service/abstract"
)

func CreateTeamHandler(ctx context.HandlerContext, req *dto.CreateTeamRequest, teamService abstract.ITeamService) (dto.CreateTeamResponse, error) {
	team, err := teamService.Create(req.Name)
	if err != nil {
		return dto.CreateTeamResponse{}, err
	}
	return dto.CreateTeamResponse{ID: team.ID, Name: team.Name}, nil
}

var _ = (*domain.Team)(nil)
