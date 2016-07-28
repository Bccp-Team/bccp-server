package scheduler

import (
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

	for {
		runId := <-sched.runRequests
		runnerId := <-sched.waitingRunners

		run, err := mysql.Db.GetRun(runId)

		if err != nil || run.Status != "waiting" {
			//FIXME the run does not exist anymore
			go sched.AddRunner(runnerId)
			continue
		}

		runner, err := mysql.Db.GetRunner(runnerId)

		if err != nil || runner.Status != "waiting" {
			//FIXME the run does not exist anymore
			go sched.AddRun(runId)
			continue
		}

		runners.StartRun(runnerId, runId)
	}
}
