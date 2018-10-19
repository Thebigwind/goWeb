package managers

import (
	//"errors"
	//"strings"
	. "github.com/bigwind/goWeb/common"
	. "github.com/bigwind/goWeb/dbservice"
)

type UserMgr struct {
	state string
}

var GlobalUserMgr *UserMgr = nil

func GetUserMgr() *UserMgr {
	return GlobalUserMgr
}

func (a *UserMgr) CheckLoginUser(userName string, userType string) error {
	db := GetDBService()
	err := db.CheckLoginDBUser(userName, userType)
	if err != nil {
		return err
	}
	return nil
}

func (a *UserMgr) CheckLoginUserPasswd(userName string, passwd string, userType string) error {
	db := GetDBService()
	err := db.CheckLoginUserPasswd(userName, passwd, userType)
	if err != nil {
		return err
	}
	return nil
}

func (a *UserMgr) GetUserId(username string) (error, int) {
	db := GetDBService()
	err, userId := db.GetUserIDByName(username)
	if err != nil {
		Logger.Errorf("GetUserIDByName failed, %s", err)
		return err, -1
	}

	return nil, userId
}
