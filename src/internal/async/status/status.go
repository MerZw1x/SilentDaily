package status

const (
	StatusQueued     = "queued"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
	StatusFailed     = "failed"
)

type CompletionStatus struct {
	UpdateID int
	Err      error
}
