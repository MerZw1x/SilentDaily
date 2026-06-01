package public

import (
	"time"

	context "silent/src/internal/context/abstract"
	"silent/src/internal/dto"
	"silent/src/internal/service/abstract"

	db "silent/src/internal/db/abstract"
)

type digestParams struct {
	TeamID int    `query:"team_id" validate:"required"`
	Date   string `query:"date"`
}

func GetDigestHandler(ctx context.HandlerContext, params digestParams, conn db.IDBConnection, digestService abstract.IDigestService) (dto.DigestResponse, error) {
	date := time.Now().AddDate(0, 0, -1)
	if params.Date != "" {
		parsed, err := time.Parse("2006-01-02", params.Date)
		if err == nil {
			date = parsed
		}
	}

	digest, err := digestService.GetByTeamAndDate(conn, params.TeamID, date)
	if err != nil {
		return dto.DigestResponse{}, err
	}
	if digest == nil {
		return dto.DigestResponse{}, notFoundError("digest not found for this team and date")
	}

	return dto.DigestResponse{
		TeamID:     digest.TeamID,
		Date:       digest.Date,
		LeadDigest: digest.LeadDigest,
	}, nil
}

type notFoundError string

func (e notFoundError) Error() string      { return string(e) }
func (e notFoundError) StatusCode() int    { return 404 }
