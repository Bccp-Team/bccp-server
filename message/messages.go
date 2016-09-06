package message

const (
	Ping = iota
	Ack
	Finish
	Kill
	Logs
	Error
	Run
)

type SubscribeRequest struct {
	Token       string
	Name        string
	Concurrency int64
}

type SubscribeAnswer struct {
	ClientUID int64
}

type ClientRequest struct {
	Kind    int64
	JobID   int64
	Logs    []string
	Message string
	Status  string
}

type ServerRequest struct {
	Kind  int64
	JobID int64
	Run   *RunRequest
}

type RunRequest struct {
	Init       string
	Repo       string
	Name       string
	UpdateTime uint64
	Timeout    uint64
}
