package runners

import (
	"encoding/gob"
	"errors"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/bccp-server/ischeduler"
	"github.com/bccp-server/message"
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
	log.Printf("WARNING: runner: %v: start connection", conn.RemoteAddr())
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	var connection message.SubscribeRequest
	err := decoder.Decode(&connection)

	if err != nil {
		log.Printf("WARNING: runner: %v: failed to decode connection: %v", conn.RemoteAddr(), err.Error())
		return
	}

	if connection.Token != *token {
		log.Printf("WARNING: runner: %v: wrong token: %v", conn.RemoteAddr(), connection.Token)
		return
	}

	uid, err := mysql.Db.AddRunner(conn.RemoteAddr().String())

	if err != nil {
		log.Printf("WARNING: runner: %v: failed to add runner: %v", conn.RemoteAddr(), err.Error())
		return
	}
	defer mysql.Db.UpdateRunner(uid, "dead")

	answer := message.SubscribeAnswer{ClientUID: uid}

	err = encoder.Encode(&answer)

	if err != nil {
		log.Printf("WARNING: runner: %v: failed to encode ack: %v", conn.RemoteAddr(), err.Error())
		return
	}

	client := &clientInfo{uid: uid, conn: conn, encoder: encoder, decoder: decoder}

	runnerMaps[uid] = client
	for i := 0; i < connection.Concurrency; i = i + 1 {
		sched.AddRunner(uid)
	}

	for {
		var clientReq message.ClientRequest
		err = decoder.Decode(&clientReq)

		if err != nil {
			log.Printf("WARNING: runner: %v: failed to decode request: %v", conn.RemoteAddr(), err.Error())
			return
		}

		switch clientReq.Kind {
		case message.Ack:
			ack(uid)
		case message.Finish:
			finish(uid, clientReq.JobId, clientReq.Status)
		case message.Logs:
			logs(uid, clientReq.JobId, clientReq.Logs)
		case message.Error:
			log.Printf("WARNING: runner: %v: receive error: %v", conn.RemoteAddr(), clientReq.Message)
		default:
			log.Printf("WARNING: runner: %v: unknow request: %v", conn.RemoteAddr(), clientReq.Kind)
			return
		}
	}
}

func KillRunner(uid int) {
	runner, ok := runnerMaps[uid]

	if !ok {
		log.Printf("WARNING: runner: %v: kill an inexistant runner", uid)
		return
	}

	runner.conn.Close()
	mysql.Db.UpdateRunner(uid, "dead")
}

func KillRun(uid, jobId int) {
	runner, ok := runnerMaps[uid]

	if !ok {
		log.Printf("WARNING: runner: kill a run on an inexistant runner (%v - %v)", uid, jobId)
		return
	}

	servReq := &message.ServerRequest{Kind: message.Kill, JobId: jobId, Run: nil}

	go func() {
		err := runner.encoder.Encode(servReq)

		if err != nil {
			log.Printf("WARNING: runner: failed to send kill request (%v - %v)", uid, jobId)
			return
		}

	}()
}

func StartRun(uid, jobId int) error {
	runner, ok := runnerMaps[uid]

	if !ok {
		log.Printf("WARNING: runner: %v: run on an inexistant runner", uid)
		return errors.New("the runner does not exist")
	}

	run, err := mysql.Db.GetRun(jobId)

	if err != nil {
		log.Printf("WARNING: runner: %v: %v", uid, err.Error())
		return err
	}

	repo, err := mysql.Db.GetRepo(run.Repo)

	if err != nil {
		log.Printf("WARNING: runner: (%v - %v): %v", uid, jobId, err.Error())
		return err
	}

	batch, err := mysql.Db.GetBatch(run.Batch)

	if err != nil {
		log.Printf("WARNING: runner: (%v - %v): %v", uid, jobId, err.Error())
		return err
	}

	runReq := &message.RunRequest{Init: batch.Init_script, Repo: repo.Ssh,
		Name: repo.Repo, UpdateTime: uint(batch.Update_time),
		Timeout: uint(batch.Timeout)}
	servReq := &message.ServerRequest{Kind: message.Run, JobId: jobId, Run: runReq}

	go func() {
		err := runner.encoder.Encode(servReq)

		if err != nil {
			log.Printf("WARNING: runner: failed to send run request (%v - %v): %v", uid, jobId, err.Error())
			mysql.Db.UpdateRunner(uid, "dead")
			runner.conn.Close()
			sched.AddRun(jobId)
			return
		}

		mysql.Db.LaunchRun(jobId, uid)
		//mysql.Db.UpdateRunner(uid, "running")
	}()

	return nil
}

func PingRunner(uid int) error {
	runner, ok := runnerMaps[uid]

	if !ok {
		log.Printf("WARNING: runner: %v: run on an inexistant runner", uid)
		return errors.New("the runner does not exist")
	}

	servReq := &message.ServerRequest{Kind: message.Ping, Run: nil}

	go func() {
		err := runner.encoder.Encode(servReq)
		if err != nil {
			log.Printf("WARNING: runner: failed to send ping request %v: %v", uid, err.Error())
			return
		}
	}()

	return nil
}

func ack(uid int) {
	r, err := mysql.Db.GetRunner(uid)
	if err != nil {
		log.Printf("WARNING: runner: ack on unknow runner %v: %v", uid, err.Error())
		return
	}
	err = mysql.Db.UpdateRunner(r.Id, r.Status)
	if err != nil {
		log.Printf("WARNING: runner: can't update runner %v: %v", uid, err.Error())
	}
}

func finish(uid int, jobId int, status string) {
	err := mysql.Db.UpdateRunStatus(jobId, status)
	if err != nil {
		log.Printf("WARNING: runner: update on unknow run %v: %v", jobId, err.Error())
	}
	err = mysql.Db.UpdateRunner(uid, "waiting")
	if err != nil {
		log.Printf("WARNING: runner: update on unknow runner %v: %v", uid, err.Error())
		return
	}
	sched.AddRunner(uid)
}

func logs(uid int, jobId int, logs []string) {
	err := mysql.Db.UpdateRunLogs(jobId, strings.Join(logs, "\n")+"\n")
	if err != nil {
		log.Printf("WARNING: runner: update on unknow run %v: %v", jobId, err.Error())
	}
}
