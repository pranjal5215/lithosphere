package main

import (
	"container/list"
	"errors"
	//"fmt"
	"sync"
)

type Worker struct {
	Id          string
	Results     chan int64
	FreeWorkers chan *Worker
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

func (wp *WorkerPool) GetWorker() (Worker, error) {
	// Get a new worker from our pool, create if required.
	wp.lk.Lock()
	defer wp.lk.Unlock()
	// Max workers already reached, so return error.
	if len(wp.totalUsedWorkers) >= wp.totalMaxWorkers {
		return Worker{}, errors.New("lithosphere:workpool: Max Workers Reached")
	} else {
		//RPOP
		e := wp.totalFreeWorkers.Back()
		if e != nil {
			//FreeWorkers List is empty.
			worker := wp.totalFreeWorkers.Remove(e)
		} else {
			//Create a new worker.
			worker := wp.createWorker()
		}
		return worker, nil
	}
}

func (wp *WorkerPool) createWorker() Worker {
	// Create a new worker.
	id := "workerid"
	return &Worker{id, manager.Results, manager.FreeWorkers}
}
