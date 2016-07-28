package runners

import (
	"encoding/gob"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/bccp-server/ischeduler"
	"github.com/bccp-server/mysql"
)

var (
	runnerMaps map[int]*clientInfo
	sched      ischeduler.IScheduler
)

func WaitRunners(isched ischeduler.IScheduler, service string, token string) {

	runnerMaps = make(map[int]*clientInfo)
	sched = isched
	log.Printf(service)

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)

	if err != nil {
		log.Panic(err)
	}

	//FIXME: TLS, toussa toussa
	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf(err.Error())
			continue
		}

		go handleClient(conn, &token)
	}
}

type clientInfo struct {
	uid        int
	currentRun uint
	conn       net.Conn
	mut        sync.Mutex
	encoder    *gob.Encoder
	decoder    *gob.Decoder
}

func handleClient(conn net.Conn, token *string) {
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	var connection SubscribeRequest
	err := decoder.Decode(&connection)

	if err != nil {
		log.Printf(err.Error())
		return
	}

	if connection.Token != *token {
		log.Printf("bad token receive: %v", connection.Token)
		return
	}

	uid, err := mysql.Db.AddRunner(conn.RemoteAddr().String())

	if err != nil {
		//FIXME error
	}

	answer := SubscribeAnswer{ClientUID: uid}

	err = encoder.Encode(&answer)

	if err != nil {
		log.Printf(err.Error())
		mysql.Db.UpdateRunner(uid, "dead")
		return
	}

	client := &clientInfo{uid: uid, conn: conn, encoder: encoder, decoder: decoder}

	runnerMaps[uid] = client
	sched.AddRunner(uid)

	for {
		var clientReq ClientRequest
		err = decoder.Decode(&clientReq)

		if err != nil {
			log.Printf(err.Error())
			return
		}

		switch clientReq.Kind {
		case Ack:
			ack(uid)
		case Finish:
			finish(uid, clientReq.JobId, clientReq.Status)
		case Logs:
			logs(uid, clientReq.JobId, clientReq.Logs)
		case Error:
		default:
		}
	}
}

//FIXME: we should take the runner and the run as parameter
func KillRun(uid int) {
	//FIXME
}

func StartRun(uid, jobId int) error {
	runner, ok := runnerMaps[uid]

	if !ok {
		//FIXME error
	}

	run, err := mysql.Db.GetRun(jobId)

	if err != nil {
		//FIXME error
	}

	//FIXME: retreive Name, Repo, UpdateTime and Timeout
	runReq := &RunRequest{Init: "ls; sleep 5; echo toto", Repo: run.Repo,
		Name: "toto", JobId: jobId, UpdateTime: 1, Timeout: 60}
	servReq := &ServerRequest{Kind: Run, Run: runReq}

	go func() {
		err := runner.encoder.Encode(servReq)

		if err != nil {
			//FIXME error
			return
		}

		mysql.Db.LaunchRun(jobId, uid)
		mysql.Db.UpdateRunner(uid, "running")
	}()

	return nil
}

func PingRunner(uid int) error {
	runner, ok := runnerMaps[uid]

	if !ok {
		//FIXME error
	}

	servReq := &ServerRequest{Kind: Ping, Run: nil}

	go func() {
		err := runner.encoder.Encode(servReq)
		if err != nil {
			//FIXME error
			return
		}
	}()

	return nil
}

func ack(uid int) {
	r, err := mysql.Db.GetRunner(uid)
	if err != nil {
		//FIXME error
	}
	err = mysql.Db.UpdateRunner(r.Id, r.Status)
	if err != nil {
		//FIXME error
	}
}

func finish(uid int, jobId int, status string) {
	err := mysql.Db.UpdateRunStatus(jobId, status)
	if err != nil {
		//FIXME error
	}
	err = mysql.Db.UpdateRunner(uid, "waiting")
	if err != nil {
		//FIXME error
	}
}

func logs(uid int, jobId int, logs []string) {
	err := mysql.Db.UpdateRunLogs(jobId, strings.Join(logs, ""))
	if err != nil {
		//FIXME error
	}
}
