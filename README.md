# lithosphere
Uppermost layer for Apps to interact with any external interface

Implement a concurrency  worker pool which everybody connects to get a worker to do a job.
Implments function closure to be called in for concurrency workers.
First implementation only to be for one among concurrency

Lithosphere  handles concurrency part and merging to be followed in next iteration.


Inspiration for lithosphere has been derived from

[golang blog.](http://blog.golang.org/concurrency-is-not-parallelism)
[github.](https://github.com/Jeffail/tunny)

and from discussions like 
[stackoverfllow.](http://stackoverflow.com/questions/23837368/idiomatic-variable-size-worker-pool-in-go)
[gobyexample.](https://gobyexample.com/worker-pools)


Functionality for Datastore connection pools and their monitoring has been moved to mantle.
Handle multiple workers so that all can be monitored when they access common connection pool:
1) Redis
2) Cassandra
3) MySQL
4) Memcache



Handle multiple workers for concurrency management.
```
package main

import (
	"fmt"
	"lithosphere"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	results := make(chan string)
	numworkers := 3
	for i := 0; i < numworkers; i++ {
		lithosphere.MainManager.ManageCoreJob(results, hello, strconv.Itoa(i))
	}
	
	for i := 0; i < numworkers; i++ {
		v := <-results
		fmt.Println(v)
	}

}

func hello(msg string) string {

	fmt.Println("HEllo ", msg)
	i := time.Duration(rand.Int31n(10000))
	fmt.Println(i)
	time.Sleep(i * time.Millisecond)

	return "Success"
}

```
