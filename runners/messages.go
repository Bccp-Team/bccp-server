package runners

const (
	Ping = iota
	Ack
	Finish
	Kill
	Logs
	Error
)

type SubscribeRequest struct {
	Token string
}

type SubscribeAnswer struct {
	ClientUID uint
}

type ClientRequest struct {
	Kind        int
	Logs        []string
	ReturnValue int
}

type ServerRequest struct {
	Kind int
	run  *RunRequest
}

type RunRequest struct {
	Init       string
	Repo       string
	Name       string
	RunId      uint
	UpdateTime uint
	Timeout    uint
}
