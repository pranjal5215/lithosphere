package lithosphere

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goibibo/gopool"
	"github.com/goibibo/t-settings"
)

const (
	MIN_CONNS = 10
	MAX_CONNS = 100
)

type MySqlPool struct {
	*pool.Pool
}

var mySqlPoolMap map[string]MySqlPool

func init() {
	mySqlPoolMap = make(map[string]MySqlPool)
}

func newMySqlPool(connStr string) MySqlPool {
	create := func() interface{} {
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			return nil
		}
		return db
	}
	p := pool.CreatePool(MIN_CONNS, MAX_CONNS,
		create,
		func(i interface{}) {
			db := i.(sql.DB)
			err := db.Close()
			if err != nil {
				log.Error("mysql db closing failed")
			}
		})
	return MySqlPool{p}
}

func GetMySqlPool(vertical string) MySqlPool {
	if pool, ok := mySqlPoolMap[vertical]; ok {
		return pool
	} else {
		configs := settings.GetConfigsFor("mysql", vertical)
		connectionString := settings.ConstructMysqlPath(configs)
		mySqlPoolMap[vertical] = newMySqlPool(connectionString)
		return mySqlPoolMap[vertical]
	}
}

func (m MySqlPool) Get() sql.DB {
	d := m.Acquire()
	dbPtr := (d).(*sql.DB)
	db := *dbPtr
	return db
}

func (m MySqlPool) Return(db sql.DB) {
	m.Release(&db)
}
