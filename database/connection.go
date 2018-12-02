package database

import (
	"io/ioutil"
	"runtime"

	"github.com/jackc/pgx"
)

// Connect to postgrss
func Connect() *pgx.ConnPool {
	runtime.GOMAXPROCS(runtime.NumCPU())
	connection := pgx.ConnConfig{
		Host:     "localhost",
		User:     "dzaytsev",
		Password: "126126",
		Database: "rsoi",
		Port:     5432,
	}

	var err error
	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: connection, MaxConnections: 50})
	if err != nil {
		panic(err)
	}

	err = createShema(db)
	if err != nil {
		panic(err)
	}

	return db
}

func createShema(db *pgx.ConnPool) error {
	sql, err := ioutil.ReadFile("database/scheme.sql")
	if err != nil {
		return err
	}
	shema := string(sql)

	_, err = db.Exec(shema)
	return err
}
