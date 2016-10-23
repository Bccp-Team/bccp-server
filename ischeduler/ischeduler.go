package ischeduler

import . "github.com/Bccp-Team/bccp-server/proto/api"

type IScheduler interface {
	AddRun(run *Run)
	AddRunner(runner int64)
	Start()
}
