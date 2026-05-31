package status

// Статусы обработки асинхронных задач
const (
	StatusQueued     = "queued"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
	StatusFailed     = "failed"
)

// CompletionStatus — результат выполнения воркера
type CompletionStatus struct{}
