package lithosphere

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/goibibo/gopool"
	"github.com/goibibo/t-settings"
)

const (
	MEMCACHE_MIN_CONNS = 10
	MEMCACHE_MAX_CONNS = 100
)

type MemCachePool struct {
	*pool.Pool
}

var memCachePoolMap map[string]MemCachePool

func init() {
	memCachePoolMap = make(map[string]MemCachePool)
}

func newMemCachePool(connStr string) MemCachePool {
	create := func() interface{} {
		db := memcache.New(connStr)
		fmt.Println(db)
		return db
	}
	p := pool.CreatePool(MEMCACHE_MIN_CONNS, MEMCACHE_MAX_CONNS,
		create,
		func(i interface{}) {

		})
	return MemCachePool{p}
}

func GetMemCachePool(vertical string) MemCachePool {
	if pool, ok := memCachePoolMap[vertical]; ok {
		return pool
	} else {
		configs := settings.GetConfigsFor("memcache", vertical)
		connectionString := settings.ConstructMysqlPath(configs)
		memCachePoolMap[vertical] = newMemCachePool(connectionString)
		return memCachePoolMap[vertical]
	}
}

func (m MemCachePool) Get() memcache.Client {
	d := m.Acquire()
	dbPtr := d.(*memcache.Client)
	db := *dbPtr
	return db
}

func (m MemCachePool) Return(db memcache.Client) {
	fmt.Println(db)
	m.Release(&db)
}
