package impl

import (
	"time"

	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	repoabstract "silent/src/internal/repository/abstract"
)

type DigestService struct {
	digestRepo repoabstract.IDigestRepository
}

func NewDigestService(digestRepo repoabstract.IDigestRepository) *DigestService {
	return &DigestService{digestRepo: digestRepo}
}

func (s *DigestService) GetByTeamAndDate(conn db.IDBConnection, teamID int, date time.Time) (*domain.Digest, error) {
	return s.digestRepo.GetByTeamIDAndDate(conn, teamID, date)
}
