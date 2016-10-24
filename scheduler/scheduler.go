package scheduler

import (
	"log"
	"reflect"

	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/Bccp-Team/bccp-server/runners"

	. "github.com/Bccp-Team/bccp-server/proto/api"
)

var (
	DefaultScheduler Scheduler
)

const (
	maxPriority = 9
)

type Scheduler struct {
	runRequests    [maxPriority + 1]chan int64
	r              [maxPriority + 1]reflect.SelectCase
	waitingRunners chan int64
}

func (sched *Scheduler) AddRun(run *Run) {
	p := run.Priority
	if p > maxPriority {
		p = maxPriority
	}
	sched.runRequests[p] <- run.Id
}

func (sched *Scheduler) AddRunner(runner int64) {
	sched.waitingRunners <- runner
}

func (sched *Scheduler) GetNextRun() int64 {
	for i := range sched.runRequests {
		select {
		case r := <-sched.runRequests[maxPriority-i]:
			return r
		default:
		}
	}

	_, r, _ := reflect.Select(sched.r[:])
	return r.Int() //r
}

func (sched *Scheduler) Start() {
	for i := 0; i <= maxPriority; i++ {
		sched.runRequests[i] = make(chan int64, 4096)
		sched.r[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(sched.runRequests[i])}
	}

	sched.waitingRunners = make(chan int64, 4096)

	dbRunners := mysql.Db.ListRunners(map[string]string{"status": "waiting"}, 0, 0)

	for _, runner := range dbRunners {
		log.Printf("scheduler: add runner %v", runner.Id)
		mysql.Db.UpdateRunner(runner.Id, "dead")
	}

	runs, err := mysql.Db.ListRuns(map[string]string{"status": "waiting"}, 0, 0)
	if err != nil {
		//FIXME
	}

	for _, run := range runs {
		log.Printf("scheduler: add run %v", run.Id)
		go sched.AddRun(run)
	}

	for {
		runID := sched.GetNextRun()
		log.Printf("scheduler: pop run %v", runID)
		runnerID := <-sched.waitingRunners
		log.Printf("scheduler: pop runner %v", runnerID)

		run, err := mysql.Db.GetRun(runID)
		if err != nil || run.Status != "waiting" {
			//FIXME the run does not exist anymore
			log.Printf("scheduler: push runner %v", runnerID)
			go sched.AddRunner(runnerID)
			continue
		}

		runner, err := mysql.Db.GetRunner(runnerID)
		if err != nil || runner.Status != "waiting" {
			//FIXME the run does not exist anymore
			log.Printf("scheduler: push run %v", runID)
			go sched.AddRun(run)
			continue
		}

		runners.StartRun(runnerID, runID)
	}
}
