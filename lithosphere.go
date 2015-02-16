package lithosphere

import (
	"fmt"
)

type jobFunc func(string) string

type Manager struct {
	CorePool  WorkerPool
	RedisWorkerPool WorkerPool
}

var MainManager Manager

func (m Manager) ManageCoreJob(results chan string, f jobFunc) {
	w, err := m.CorePool.GetWorker()
	if err != nil {
		results <- ""
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		defer m.CorePool.returnWorker(w)
		w.doJob(results, f)
	}()
}


func (m Manager) ManageRedisJob(results chan string, f jobFunc) {
	w, err := m.RedisWorkerPool.GetWorker()
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
		w.doJob(results, f)
	}
}
