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
	Concurrency int
}

type SubscribeAnswer struct {
	ClientUID int
}

type ClientRequest struct {
	Kind    int
	JobID   int
	Logs    []string
	Message string
	Status  string
}

type ServerRequest struct {
	Kind  int
	JobID int
	Run   *RunRequest
}

type RunRequest struct {
	Init       string
	Repo       string
	Name       string
	UpdateTime uint
	Timeout    uint
}
