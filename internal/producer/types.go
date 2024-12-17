package producer

type toggleRequest struct {
	WorkerStatus bool `json:"workerStatus"`
}
