package dispatcher

import (
	"context"
	"log"
	"time"

	"silent/src/internal/async/semaphore"
	"silent/src/internal/async/status"
	workerconfig "silent/src/internal/async/worker/config"
	extractworker "silent/src/internal/async/worker/extract"
	db "silent/src/internal/db/abstract"
	repoabstract "silent/src/internal/repository/abstract"
)

type Dispatcher struct {
	ch          chan status.CompletionStatus
	sem         semaphore.Semaphore
	conn        db.IDBConnection
	updateRepo  repoabstract.IDailyUpdateRepository
	workerCfg   workerconfig.WorkerConfig
	sleepTime   time.Duration
	maxAttempts int
}

func NewDispatcher(
	conn db.IDBConnection,
	updateRepo repoabstract.IDailyUpdateRepository,
	workerCfg workerconfig.WorkerConfig,
	sleepTime time.Duration,
	maxAttempts int,
	maxWorkers int,
) *Dispatcher {
	return &Dispatcher{
		ch:          make(chan status.CompletionStatus, maxWorkers*2),
		sem:         semaphore.NewSemaphore(maxWorkers),
		conn:        conn,
		updateRepo:  updateRepo,
		workerCfg:   workerCfg,
		sleepTime:   sleepTime,
		maxAttempts: maxAttempts,
	}
}

func (d *Dispatcher) RunWorkers(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("[dispatcher] RunWorkers stopped: %s", ctx.Err())
			return
		default:
			d.sem.Acquire()
			d.fetchAndRun(ctx)
		}
	}
}

func (d *Dispatcher) ProcessCompletion(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("[dispatcher] ProcessCompletion stopped: %s", ctx.Err())
			return
		case result := <-d.ch:
			if result.Err != nil {
				log.Printf("[dispatcher] worker error for update %d: %v", result.UpdateID, result.Err)
				update, err := d.updateRepo.GetOneQueued(d.conn)
				if err != nil || update == nil {
					continue
				}
				if update.Attempts >= d.maxAttempts {
					_ = d.updateRepo.SetStatus(d.conn, result.UpdateID, status.StatusFailed)
				} else {
					_ = d.updateRepo.SetStatusAndIncrementAttempts(d.conn, result.UpdateID, status.StatusQueued)
				}
			} else {
				_ = d.updateRepo.SetStatus(d.conn, result.UpdateID, status.StatusDone)
			}
		default:
			time.Sleep(d.sleepTime)
		}
	}
}

func (d *Dispatcher) fetchAndRun(ctx context.Context) {
	txConn := d.conn.BeginTx()

	defer func() {
		if r := recover(); r != nil {
			txConn.Rollback()
			d.sem.Release()
		}
	}()

	update, err := d.updateRepo.GetOneQueued(txConn)
	if err != nil || update == nil {
		txConn.Rollback()
		d.sem.Release()
		time.Sleep(d.sleepTime)
		return
	}

	if err = d.updateRepo.SetStatusAndIncrementAttempts(txConn, update.ID, status.StatusInProgress); err != nil {
		txConn.Rollback()
		d.sem.Release()
		return
	}

	if err = txConn.Commit(); err != nil {
		d.sem.Release()
		return
	}

	go func() {
		defer d.sem.Release()
		extractworker.ExtractWorker(ctx, d.ch, d.workerCfg, update.ID)
	}()
}
