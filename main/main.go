package main

import (
	"sync"

	"github.com/bccp/api"
)

func main() {
	var wait sync.WaitGroup
	api.SetupRestAPI(&wait)
	wait.Wait()
}
