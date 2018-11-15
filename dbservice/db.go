package dbservice

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	. "github.com/bigwind/goWeb/common"
	_ "github.com/lib/pq"
)

func (db *dbService) CheckLoginDBUser(userName string, userType string) error {
	querySql := "SELECT user_id FROM user_main WHERE user_name=$1 and user_type = $2 and user_status='active';"
	rows, err := db.Query(querySql, userName, userType)
	if err != nil {
		db.QueryErr(err)
		DBLogger.Errorf("DBService: query error %s\n", err.Error())
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count += 1
	}

	if count == 0 {
		err := errors.New(fmt.Sprintf("Account %s not exist or not active!", userName))
		return err
	}
	return nil
}

func (db *dbService) CheckLoginUserPasswd(userName string, passwd string, userType string) error {
	querySql := "SELECT user_id FROM user_main WHERE user_name=$1 AND password = $2 AND user_type = $3 AND user_status='active';"
	fmt.Println("xoxo")
	rows, err := db.Query(querySql, userName, passwd, userType)
	if err != nil {
		db.QueryErr(err)
		DBLogger.Errorf("DBService: query error %s\n", err.Error())
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count += 1
	}

	if count == 0 {
		err := errors.New(fmt.Sprintf("passwd not right or not active!"))
		return err
	}
	return nil
}

func (db *dbService) GetUserIDByName(userName string) (error, int) {
	DBLogger.Infof("start to query userName:", userName)
	querySql := "SELECT user_id FROM user_main WHERE user_name=$1;"
	rows, err := db.Query(querySql, userName)
	if err != nil {
		db.QueryErr(err)
		DBLogger.Errorf("DBService: query error %s\n", err.Error())
		return err, -1
	}
	defer rows.Close()

	var userID int = -1
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return err, -1
		}
	}

	return nil, userID
}

func (db *dbService) BatchMode1(id int, addressTags string) error {
	tx, err := db.db.Begin()
	defer func() {
		if tx != nil {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()

	insertSql := ""
	tags := strings.Split(addressTags, ",")
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		insertSql += "INSERT INTO batch_test values(" + strconv.Itoa(id) + ");"
	}
	if insertSql == "" {
		return nil
	}
	Logger.Infoln("insertSql=%s", insertSql)

	err = db.TxExec(tx, insertSql)
	if err != nil {
		Logger.Errorf(err.Error())
		err = errors.New("insert address tags failed.")
		return err
	}
	return nil
}
