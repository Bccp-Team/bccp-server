package runners

import (
	"encoding/gob"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Bccp-Team/bccp-server/ischeduler"
	"github.com/Bccp-Team/bccp-server/message"
	"github.com/Bccp-Team/bccp-server/mysql"
)

var (
	runnerMaps map[int64]*clientInfo
	sched      ischeduler.IScheduler
	mutex      sync.RWMutex
)

func WaitRunners(isched ischeduler.IScheduler, service string, token string) {
	runnerMaps = make(map[int64]*clientInfo)
	sched = isched
	mutex := sync.RWMutex{}
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
	uid         int64
	currentRun  uint64
	conn        net.Conn
	mut         sync.Mutex
	encoder     *gob.Encoder
	decoder     *gob.Decoder
	pingChannel chan bool
}

func cleanupClient(uid int64) {
	runnerId := strconv.FormatInt(uid, 10)
	runs, _ := mysql.Db.ListRuns(map[string]string{"runner": runnerId,
		"status": "running"}, 0, 0)
	for _, run := range runs {
		log.Printf("runner: kill %v", run.RunnerId)
		mysql.Db.UpdateRunStatus(run.Id, "killed")
		id, err := mysql.Db.AddRun(run.RepoId, run.Batch, run.Priority)
		if err != nil {
			log.Printf("runner: could not reschedul %v: %v", run.RunnerId, err.Error())
			continue
		}
		nrun, err := mysql.Db.GetRun(id)
		if err != nil {
			log.Printf("runner: could not reschedul %v: %v", run.RunnerId, err.Error())
			continue
		}
		sched.AddRun(nrun)
	}
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

	uid, err := mysql.Db.AddRunner(conn.RemoteAddr().String(), connection.Name)
	if err != nil {
		log.Printf("WARNING: runner: %v: failed to add runner: %v", conn.RemoteAddr(), err.Error())
		return
	}

	defer cleanupClient(uid)
	defer mysql.Db.UpdateRunner(uid, "dead")

	answer := message.SubscribeAnswer{ClientUID: uid}

	err = encoder.Encode(&answer)
	if err != nil {
		log.Printf("WARNING: runner: %v: failed to encode ack: %v", conn.RemoteAddr(), err.Error())
		return
	}

	client := &clientInfo{uid: uid, conn: conn, encoder: encoder, decoder: decoder, pingChannel: make(chan bool)}

	mutex.Lock()
	runnerMaps[uid] = client
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(runnerMaps, uid)
		mutex.Unlock()
	}()

	for i := int64(0); i < connection.Concurrency; i = i + 1 {
		sched.AddRunner(uid)
	}

	go client.ping()

	for {
		var clientReq message.ClientRequest
		err = decoder.Decode(&clientReq)
		if err != nil {
			log.Printf("WARNING: runner: %v: failed to decode request: %v", conn.RemoteAddr(), err.Error())
			return
		}

		switch clientReq.Kind {
		case message.Ack:
			client.ack()
		case message.Finish:
			client.finish(clientReq.JobID, clientReq.Status)
		case message.Logs:
			client.logs(clientReq.JobID, clientReq.Logs)
		case message.Error:
			log.Printf("WARNING: runner: %v: receive error: %v", conn.RemoteAddr(), clientReq.Message)
		default:
			log.Printf("WARNING: runner: %v: unknow request: %v", conn.RemoteAddr(), clientReq.Kind)
			return
		}
	}
}

func KillRunner(uid int64) {
	mutex.RLock()
	runner, ok := runnerMaps[uid]
	mutex.RUnlock()
	if !ok {
		log.Printf("WARNING: runner: %v: kill an inexistant runner", uid)
		return
	}

	runner.conn.Close()
	mysql.Db.UpdateRunner(uid, "dead")
}

func KillRun(uid, jobID int64) {
	mutex.RLock()
	runner, ok := runnerMaps[uid]
	mutex.RUnlock()
	if !ok {
		log.Printf("WARNING: runner: kill a run on an inexistant runner (%v - %v)", uid, jobID)
		return
	}

	servReq := &message.ServerRequest{Kind: message.Kill, JobID: jobID, Run: nil}

	go func() {
		err := runner.encoder.Encode(servReq)
		if err != nil {
			log.Printf("WARNING: runner: failed to send kill request (%v - %v)", uid, jobID)
			return
		}
	}()
}

func StartRun(uid, jobID int64) error {
	mutex.RLock()
	runner, ok := runnerMaps[uid]
	mutex.RUnlock()
	if !ok {
		log.Printf("WARNING: runner: %v: run on an inexistant runner", uid)
		return errors.New("the runner does not exist")
	}

	run, err := mysql.Db.GetRun(jobID)
	if err != nil {
		log.Printf("WARNING: runner: %v: %v", uid, err.Error())
		return err
	}

	repo, err := mysql.Db.GetRepo(run.RepoId)
	if err != nil {
		log.Printf("WARNING: runner: (%v - %v): %v", uid, jobID, err.Error())
		return err
	}

	batch, err := mysql.Db.GetBatch(run.Batch)
	if err != nil {
		log.Printf("WARNING: runner: (%v - %v): %v", uid, jobID, err.Error())
		return err
	}

	runReq := &message.RunRequest{Init: batch.InitScript, Repo: repo.Ssh,
		Name: repo.Repo + "_" + strconv.FormatInt(jobID, 10), UpdateTime: uint64(batch.UpdateTime),
		Timeout: uint64(batch.Timeout)}
	servReq := &message.ServerRequest{Kind: message.Run, JobID: jobID, Run: runReq}

	go func() {
		mysql.Db.LaunchRun(jobID, uid)
		err := runner.encoder.Encode(servReq)
		if err != nil {
			log.Printf("WARNING: runner: failed to send run request (%v - %v): %v", uid, jobID, err.Error())
			mysql.Db.UpdateRunner(uid, "dead")
			runner.conn.Close()
			sched.AddRun(run)
			return
		}
	}()

	return nil
}

func PingRunner(uid int64) error {
	mutex.RLock()
	runner, ok := runnerMaps[uid]
	mutex.RUnlock()
	if !ok {
		log.Printf("WARNING: runner: %v: ping on an inexistant runner", uid)
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

func (client *clientInfo) ack() {
	client.pingChannel <- true
	r, err := mysql.Db.GetRunner(client.uid)
	if err != nil {
		log.Printf("WARNING: runner: ack on unknow runner %v: %v", client.uid, err.Error())
		return
	}

	err = mysql.Db.UpdateRunner(r.Id, r.Status)
	if err != nil {
		log.Printf("WARNING: runner: can't update runner %v: %v", client.uid, err.Error())
	}
}

func (client *clientInfo) ping() {
	//FIXME fetch info from config
	timer := time.After(time.Minute)
	tick := time.Tick(5 * time.Second)

	for {
		select {
		case <-client.pingChannel:
			timer = time.After(time.Minute)
		case <-tick:
			if PingRunner(client.uid) != nil {
				return
			}
		case <-timer:
			log.Printf("WARNING: runner: %v: timeout", client.conn.RemoteAddr())
			client.conn.Close()
			return
		}
	}
}

func (client *clientInfo) finish(jobID int64, status string) {
	run, err := mysql.Db.GetRun(jobID)
	if err != nil {
		log.Printf("WARNING: runner: update on unknow run %v - %v: %v", client.uid, jobID, err.Error())
		KillRunner(client.uid)
		return
	}

	if run.RunnerId != int64(client.uid) {
		log.Printf("WARNING: runner: runner update wrong run %v %v: %v", client.uid, jobID, err.Error())
		KillRunner(client.uid)
		return
	}

	if run.Status == "running" {
		err = mysql.Db.UpdateRunStatus(jobID, status)
		if err != nil {
			log.Printf("WARNING: runner: update on unknow run %v: %v", jobID, err.Error())
		}
	}

	err = mysql.Db.UpdateRunner(client.uid, "waiting")
	if err != nil {
		log.Printf("WARNING: runner: update on unknow runner %v: %v", client.uid, err.Error())
		KillRunner(client.uid)
		return
	}

	sched.AddRunner(client.uid)
}

func (client *clientInfo) logs(jobID int64, logs []string) {
	run, err := mysql.Db.GetRun(jobID)
	if err != nil {
		log.Printf("WARNING: runner: update on unknow run %v - %v: %v", client.uid, jobID, err.Error())
		KillRunner(client.uid)
		return
	}

	if run.RunnerId != int64(client.uid) {
		log.Printf("WARNING: runner: runner update wrong run %v: %v", client.uid, jobID)
		KillRunner(client.uid)
		return
	}

	if run.Status != "running" {
		//FIXME logs
		log.Printf("WARNING: runner: update on finished run %v: %v", client.uid, jobID)
		KillRun(client.uid, jobID)
		return
	}

	if logs != nil || len(logs) == 0 {
		return
	}

	err = mysql.Db.UpdateRunLogs(jobID, strings.Join(logs, "\n")+"\n")
	if err != nil {
		log.Printf("WARNING: runner: update on unknow run %v: %v", jobID, err.Error())
	}
}
