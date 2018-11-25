package database

import (
	"io/ioutil"
	"runtime"

	"github.com/jackc/pgx"
)

// DB instance
var DB *pgx.ConnPool

// Connect to postgrss
func Connect() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	connection := pgx.ConnConfig{
		Host:     "localhost",
		User:     "dzaytsev",
		Password: "126126",
		Database: "rsoi",
		Port:     5432,
	}

	var err error
	DB, err = pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: connection, MaxConnections: 50})
	if err != nil {
		panic(err)
	}

	err = createShema()
	if err != nil {
		panic(err)
	}
}

func createShema() error {
	sql, err := ioutil.ReadFile("scheme.sql")
	if err != nil {
		return err
	}
	shema := string(sql)

	_, err = DB.Exec(shema)
	return err
}
