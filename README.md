# lithosphere
Uppermost layer for Apps to interact with any external interface

Implement a worker pool which everybody connects to to get a worker to do a job.
Implments function closure to be called in for concurrency/API/DB workers.
First implementation only to be for one among concurrency/API/DB


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
	for i := 1; i < 3; i++ {
		lithosphere.MainManager.ManageCoreJob(results, hello, strconv.Itoa(i))
	}
	v := <-results
	fmt.Println(v)
	v = <-results
	fmt.Println(v)

}

func hello(na string) string {

	fmt.Println("HEllo ", na)
	i := time.Duration(rand.Int31n(10000))
	fmt.Println(i)
	time.Sleep(i * time.Millisecond)

	return "Success"
}

```
