package lithosphere

import (
	"fmt"
)

var MAXWORKER int = 50

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
	w, err := m.CorePool.getWorker()
	if err != nil {
		fmt.Println(err)
		go func() {
			results <- "err"
		}()
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		defer m.CorePool.returnWorker(w)
		w.doJob(results, f, inp)
	}()
}

func (m Manager) ManageRedisJob(results chan string, f jobFunc, inp string) {
	w, err := m.RedisWorkerPool.getWorker()
	if err != nil {
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		defer m.CorePool.returnWorker(w)
		w.doJob(results, f, inp)
	}()
}
