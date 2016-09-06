package scheduler

import (
	"log"

	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/Bccp-Team/bccp-server/runners"
)

var (
	DefaultScheduler Scheduler
)

type Scheduler struct {
	runRequests    chan int64
	waitingRunners chan int64
}

func (sched *Scheduler) AddRun(run int64) {
	sched.runRequests <- run
}

func (sched *Scheduler) AddRunner(runner int64) {
	sched.waitingRunners <- runner
}

func (sched *Scheduler) Start() {
	sched.runRequests = make(chan int64, 4096)
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
		go sched.AddRun(run.Id)
	}

	for {
		runID := <-sched.runRequests
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
			go sched.AddRun(runID)
			continue
		}

		runners.StartRun(runnerID, runID)
	}
}
