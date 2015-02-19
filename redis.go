package lithosphere

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/goibibo/t-settings"
	"time"
)

type RedisPool struct {
	*redis.Pool
}

var PoolMap map[string]*RedisPool

func init() {
	PoolMap = make(map[string]*RedisPool)
}

func newPool(server, password string) *RedisPool {
	r := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &RedisPool{r}
}

func GetPool(vertical string) *RedisPool {
	if pool, ok := PoolMap[vertical]; ok {
		return pool
	} else {
		configs := settings.GetConfigsFor("redis", vertical)
		connectionString := settings.ConstructRedisPath(configs)
		PoolMap[vertical] = newPool(connectionString, "")
		return PoolMap[vertical]
	}
}

func GetAllActiveConnections() string {
	var msg string
	for ver, pool := range PoolMap {
		msg += fmt.Sprintln("Vertical", ver, " : ", pool.ActiveCount())
	}
	return msg
}
