package webFile

import (
	"errors"
	"fmt"

	. "github.com/bigwind/goWeb/common"
	"github.com/bigwind/goWeb/dbservice"
	. "github.com/bigwind/goWeb/server"
)

const (
	DEFAULT_REST_SERV string = "127.0.0.1"
	DEFAULT_REST_PORT string = "9090"
)

func ServerStart() error {
	/*
	 * Step 1: initialize the logger facility.
	 */
	loggerConfig := LoggerConfig{
		Logfile:  "goWeb.log", ///var/log/goWeb/goWeb.log
		LogLevel: 0,
	}
	LoggerInit(&loggerConfig)

	/*
	 * Step 2: for local test, connect any db without read etcd.
	 */
	dbConfig := DBServiceConfig{
		Server:   "127.0.0.1",
		Port:     "5432",
		User:     "postgres",
		Password: "123456",
		Driver:   "postgres",
		DBName:   "testdb",
	}

	dbService := dbservice.NewDBService(&dbConfig)
	if dbService == nil {
		Logger.Errorf("Fail to init DB service")
		return errors.New("Fail init DB service")
	} else {
		Logger.Infof("datase connect now...")
	}

	restAddr := fmt.Sprintf(":%s", DEFAULT_REST_PORT)
	serv := NewRESTServer(restAddr)
	serv.StartRESTServer()
	return nil
}
