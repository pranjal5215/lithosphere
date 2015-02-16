package lithosphere

import (
	"code.google.com/p/go-uuid/uuid"
	"container/list"
	"errors"
	"sync"
)

type Worker struct {
	Id string
}

type WorkerPool struct {
	//Singleton object on worker pool to be created which will be shared by all.
	//PoolWorker is of type Worker
	PoolWorker       Worker
	lk               sync.Mutex
	totalMaxWorkers  int               //Maximum allowed workers.
	totalUsedWorkers map[string]Worker //Workers which are in use currently.
	totalFreeWorkers list.List         //Workers available to be picked up.
}

func (wp *WorkerPool) getWorker() (Worker, error) {
	// Get a new worker from our pool, create if required.
	wp.lk.Lock()
	defer wp.lk.Unlock()
	// Max workers already reached, so return error.
	if len(wp.totalUsedWorkers) >= wp.totalMaxWorkers {
		return Worker{}, errors.New("lithosphere:workpool: Max Workers Reached")
	} else {
		//RPOP
		e := wp.totalFreeWorkers.Back()
		var worker Worker
		if e != nil {
			//FreeWorkers List is not empty.
			wp.totalFreeWorkers.Remove(e)
			worker = e.Value.(Worker)
		} else {
			//Create a new worker.
			worker = wp.createWorker()
		}
		// Put worker in Used queue.
		// Manager should always start working
		// as soon as getWorker returns a worker.
		wp.totalUsedWorkers[worker.Id] = worker
		return worker, nil
	}
}

func (w Worker) doJob(results chan string, f jobFunc) {
	result := f("a")
	results <- result
}

func (wp *WorkerPool) createWorker() Worker {
	// Create a new worker.
	id := uuid.New()
	return Worker{id}
}

func (wp *WorkerPool) returnWorker(w Worker) {
	//Return worker from Used map tp Free queue
	wp.lk.Lock()
	defer wp.lk.Unlock()

	first := wp.totalFreeWorkers.Front()
	wp.totalFreeWorkers.InsertBefore(w, first)
	delete(wp.totalUsedWorkers, w.Id)
}
