package ischeduler

type IScheduler interface {
	AddRun(run int)
	AddRunner(runner int)
	Start()
}
