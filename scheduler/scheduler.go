package scheduler

import (
	"log"

	"github.com/bccp-server/mysql"
	"github.com/bccp-server/runners"
)

var (
	DefaultScheduler Scheduler
)

type Scheduler struct {
	runRequests    chan int
	waitingRunners chan int
}

func (sched *Scheduler) AddRun(run int) {
	sched.runRequests <- run
}

func (sched *Scheduler) AddRunner(runner int) {
	sched.waitingRunners <- runner
}

func (sched *Scheduler) Start() {
	sched.runRequests = make(chan int, 4096)
	sched.waitingRunners = make(chan int, 4096)

	runners_ := mysql.Db.ListRunners(map[string]string{"status": "waiting"})

	for _, runner := range runners_ {
		log.Printf("scheduler: add runner %v", runner.Id)
		mysql.Db.UpdateRunner(runner.Id, "dead")
	}

	runs, err := mysql.Db.ListRuns(map[string]string{"status": "waiting"})

	if err != nil {
		//FIXME
	}

	for _, run := range runs {
		log.Printf("scheduler: add run %v", run.Id)
		go sched.AddRun(run.Id)
	}

	for {
		runId := <-sched.runRequests
		log.Printf("scheduler: pop run %v", runId)
		runnerId := <-sched.waitingRunners
		log.Printf("scheduler: pop runner %v", runnerId)

		run, err := mysql.Db.GetRun(runId)

		if err != nil || run.Status != "waiting" {
			//FIXME the run does not exist anymore
			log.Printf("scheduler: push runner %v", runnerId)
			go sched.AddRunner(runnerId)
			continue
		}

		runner, err := mysql.Db.GetRunner(runnerId)

		if err != nil || runner.Status != "waiting" {
			//FIXME the run does not exist anymore
			log.Printf("scheduler: push run %v", runId)
			go sched.AddRun(runId)
			continue
		}

		runners.StartRun(runnerId, runId)
	}
}
