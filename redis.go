package lithosphere

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/goibibo/t-settings"
	"time"
)

type RedisPool struct {
	*redis.Pool
}

var redisPoolMap map[string]*RedisPool

func init() {
	redisPoolMap = make(map[string]*RedisPool)
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
	if pool, ok := redisPoolMap[vertical]; ok {
		return pool
	} else {
		configs := settings.GetConfigsFor("redis", vertical)
		connectionString := settings.ConstructRedisPath(configs)
		redisPoolMap[vertical] = newPool(connectionString, "")
		return redisPoolMap[vertical]
	}
}

func GetVerticalActiveConnections(ver string) (result int, err error) {
	if pool, ok := redisPoolMap[ver]; ok {
		return pool.ActiveCount(), nil
	} else {
		return 0, errors.New("Vertical not found in pool")
	}
}

func GetAllActiveConnections() string {
	var msg string
	for ver, pool := range redisPoolMap {
		msg += fmt.Sprintln("Vertical", ver, " : ", pool.ActiveCount())
	}
	return msg
}
