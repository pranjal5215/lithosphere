# lithosphere
Uppermost layer for Apps to interact with any external interface

Implement a worker pool which everybody connects to to get a worker to do a job.
Implments function closure to be called in for concurrency/API/DB workers.
First implementation only to be for one among concurrency/API/DB

TODO:Handle multiple workers so that all can be monitored when they access common connection pool
1) Redis
2) Cassandra
3) MySQL
4) Memcache

AND


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


Testing Memecache
```
package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/goibibo/lithosphere"
	"github.com/goibibo/t-settings"
)

func main() {
	dir := "dev"
	path := fmt.Sprintf("thor/src/settings/%s/config.json", dir)
	settings.Configure(path)

	pl := lithosphere.GetMemCachePool("flight")

	for i := 0; i <= 120; i++ {
		c, _ := pl.Get()
		fmt.Println("just got", c)
		c.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})
		it := c.Get("foo")

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(it)
		fmt.Println("returning ", c)
		pl.Return(c)
	}
}
```

Mysql Example
```
package main

import (
        "github.com/goibibo/lithosphere"
        "fmt"
        "github.com/goibibo/t-settings"

)

func main() {
        dir := "dev"
        path := fmt.Sprintf("thor/src/settings/%s/config.json", dir)
        settings.Configure(path)

        pl := lithosphere.GetMySqlPool("flight")

        for i := 0; i <= 120; i++ {
                c := pl.Get()
                fmt.Println("just got", c)

                t, _ := c.Exec("show tables")
                fmt.Println(t)

                fmt.Println("returning ", c)
                pl.Return(c)
        }
}
```