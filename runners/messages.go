package runners

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
	Token string
}

type SubscribeAnswer struct {
	ClientUID int
}

type ClientRequest struct {
	Kind    int
	JobId   int
	Logs    []string
	Message string
	Status  string
}

type ServerRequest struct {
	Kind int
	Run  *RunRequest
}

type RunRequest struct {
	Init       string
	Repo       string
	Name       string
	JobId      uint
	UpdateTime uint
	Timeout    uint
}
