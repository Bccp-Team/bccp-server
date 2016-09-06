package ischeduler

type IScheduler interface {
	AddRun(run int64)
	AddRunner(runner int64)
	Start()
}
