package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/JorniZ/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err.Error())
	}

	fmt.Println(config.DBDriver)

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the db:", err.Error())
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
