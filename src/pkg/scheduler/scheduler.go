package scheduler

import (
	"context"
	"log"

	workerconfig "silent/src/internal/async/worker/config"
	digestworker "silent/src/internal/async/worker/digest"

	"github.com/robfig/cron/v3"
)

type DigestScheduler struct {
	cron          *cron.Cron
	cfg           workerconfig.WorkerConfig
	teamIDs       []int
	memberChatIDs map[int]int64
}

func NewDigestScheduler(cfg workerconfig.WorkerConfig, teamIDs []int, memberChatIDs map[int]int64) *DigestScheduler {
	return &DigestScheduler{
		cron:          cron.New(),
		cfg:           cfg,
		teamIDs:       teamIDs,
		memberChatIDs: memberChatIDs,
	}
}

func (s *DigestScheduler) Start() {
	_, err := s.cron.AddFunc("0 8 * * *", func() {
		log.Println("[scheduler] running morning digest")
		for _, teamID := range s.teamIDs {
			digestworker.DigestWorker(context.Background(), s.cfg, teamID, s.memberChatIDs)
		}
	})
	if err != nil {
		log.Fatalf("[scheduler] failed to register cron job: %v", err)
	}
	s.cron.Start()
	log.Println("[scheduler] digest scheduler started (08:00 daily)")
}

func (s *DigestScheduler) Stop() {
	s.cron.Stop()
}
