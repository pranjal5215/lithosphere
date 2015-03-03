package lithosphere

import (
	"fmt"
)

var MAXWORKER int = 100000

type jobFunc func(string) string

type Manager struct {
	CorePool        WorkerPool
	RedisWorkerPool WorkerPool
}

var MainManager Manager

func init() {
	MainManager.CorePool.totalMaxWorkers = MAXWORKER
	mp := make(map[string]Worker)
	MainManager.CorePool.totalUsedWorkers = mp
}

func (m Manager) ManageCoreJob(results chan string, f jobFunc, inp string) {
	// ManageCoreJob is a manager to perform workers action through worker pool
	// Get a Worker
	w, err := m.CorePool.getWorker()
	if err != nil {
		// Handle error from a worker
		// TODO::Implement gotWorkerBlocking which blocks when MAXWORKERS is reached.
		fmt.Println(err)
		go func() {
			results <- "err"
		}()
		return
	}

	go func() {
		// Recovery
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		// Upon completion return worker to pool after work is done.
		defer m.CorePool.returnWorker(w)
		// Do work.
		w.doJob(results, f, inp)
	}()
}
