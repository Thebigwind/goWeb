package dbservice

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync/atomic"
	"time"

	"regexp"

	. "github.com/bigwind/goWeb/common"
	_ "github.com/lib/pq"
)

const MAX_DB_COUNT int = 10000

type dbStat struct {
	lastErr     error
	lastErrTime time.Time
	lastErrType string
	queryErr    int64
	execErr     int64
}

type dbService struct {
	server string
	port   string
	user   string
	pass   string
	dbname string
	driver string
	db     *sql.DB
	stat   *dbStat
}

func GetDBService() *dbService {
	return globalDBService
}

var globalDBService *dbService = nil

var dbReconnectPeriod = time.Second * 5
var dbReconnectCount = 10000

func NewDBService(config *DBServiceConfig) *dbService {
	stat := &dbStat{
		queryErr: 0,
		execErr:  0,
	}

	dbService := &dbService{
		server: config.Server,
		port:   config.Port,
		user:   config.User,
		pass:   config.Password,
		dbname: config.DBName,
		driver: config.Driver,
		stat:   stat,
	}

	dbOpts := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbService.user,
		dbService.pass, dbService.dbname)

	DBLogger.Errorf("%v", config)

	if dbService.server != "" {
		dbOpts += " host=" + dbService.server
	}
	if dbService.port != "" {
		dbOpts += " port=" + dbService.port
	}
	db, err := sql.Open(dbService.driver, dbOpts)
	if err != nil {
		DBLogger.Errorf("Fail to open database, error %s \n", err.Error())
		return nil
	}
	dbService.db = db

	/*
	 * DB service blocking retry dbReconnectCount * dbReconnectPeriod
	 * until connect successfully or give up.
	 */
	connected := false
	for i := 0; i <= dbReconnectCount; i++ {
		err = db.Ping()
		if err != nil {
			DBLogger.Errorf("database connection failed, retry...")
			time.Sleep(dbReconnectPeriod)
		} else {
			DBLogger.Infof("database connection successfully! \n")
			connected = true

			break
		}
	}

	if connected == false {
		DBLogger.Infof("give up to connect database, failed...! \n")
	}

	globalDBService = dbService
	return dbService
}

//add for transaction
func (db *dbService) NewDbTx() (*sql.Tx, error) {
	return db.db.Begin()
}

func (db *dbService) ExecErr(err error) {
	db.stat.lastErrTime = time.Now()
	db.stat.lastErrType = "Exec"
	db.stat.lastErr = err
	atomic.AddInt64(&db.stat.execErr, 1)
}

func (db *dbService) QueryErr(err error) {
	db.stat.lastErrTime = time.Now()
	db.stat.lastErrType = "Query"
	db.stat.lastErr = err
	atomic.AddInt64(&db.stat.queryErr, 1)
}

func (db *dbService) Query(query string, args ...interface{}) (*sql.Rows, error) {
	DBLogger.Infof("Query args:%", args)

	var err error
	var rows *sql.Rows
	for i := 0; i <= dbReconnectCount; i++ {
		rows, err = db.db.Query(query, args...)
		if err != nil {
			db.QueryErr(err)
			if IsConnectionError(err) != true {
				return nil, err
			} else {
				DBLogger.Errorf("database connection failed, retry...")
				time.Sleep(dbReconnectPeriod)
			}
		} else {
			break
		}
	}
	return rows, err
}

func (db *dbService) Exec(sqlStr string, args ...interface{}) error {
	stmt, err := db.db.Prepare(sqlStr)
	if err != nil {
		DBLogger.Errorf("DBService: prepare error %s\n", err.Error())
		return err
	}

	defer stmt.Close()

	for i := 0; i <= dbReconnectCount; i++ {
		_, err = stmt.Exec(args...)
		if err != nil {
			db.ExecErr(err)
			if IsConnectionError(err) != true {
				break
			} else {
				DBLogger.Errorf("database connection failed, retry...")
				time.Sleep(dbReconnectPeriod)
			}
		} else {
			break
		}
	}
	return err
}

func (db *dbService) TxExec(tx *sql.Tx, sqlStr string, args ...interface{}) error {
	var err error = nil

	for i := 0; i <= dbReconnectCount; i++ {
		_, err = tx.Exec(sqlStr, args...)
		if err != nil {
			db.ExecErr(err)
			if IsConnectionError(err) != true {
				break
			} else {
				DBLogger.Errorf("database connection failed, retry...")
				time.Sleep(dbReconnectPeriod)
			}
		} else {
			break
		}
	}
	return err
}

/*
	goWeb common functions
*/
func getStampTimeString(timeStamp int) string {
	//1495358104 int to "2017-05-21 17:15:04"
	int64TimeStamp := int64(timeStamp)
	stampTimeString := time.Unix(int64TimeStamp, 0).Format("2006-01-02 15:04:05")

	//DBLogger.Infof("stampTimeString:", stampTimeString)
	return stampTimeString
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func IsConnectionError(err error) bool {
	reg := regexp.MustCompile("connection?")

	return reg.FindStringIndex(err.Error()) != nil
}
