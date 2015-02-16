package lithosphere

import (
	"fmt"
)

type Manager struct {
	CorePool  WorkerPool
}

var MainManager Manager

func (m Manager) ManageCoreJob(results chan string, funcName string) {
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
		defer m.CorePool.ReturnWorker(w)
		w.doJob(results, funcName)
	}
}