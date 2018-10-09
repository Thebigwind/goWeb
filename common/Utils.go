package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

const XT_API_RET_OK string = "OK"
const XT_API_RET_ERROR string = "ERROR"
const XT_API_RET_UNKNOWN string = "UNKNOWN"

type ReturnJson struct {
	Status string      `json:"status"`
	Errmsg string      `json:"errmsg"`
	Result interface{} `json:"result"`
}

func (this *ReturnJson) AjaxReturnJson(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.WriteHeader(http.StatusAccepted)

	return json.NewEncoder(w).Encode(this)
}

func CheckReadFile(pathStr string) error {
	info, err := os.Stat(pathStr)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("No such file.")
		} else if os.IsPermission(err) {
			return errors.New("Permission denied !")
		} else {
			return err
		}
	}
	if info.IsDir() {
		return errors.New("Not a file.")
	}
	if info.Mode().Perm()&os.FileMode(128) == 0 {
		return errors.New("Cannot read the file. Permission denied !")
	}
	return nil
}

func FileExists(pathStr string) error {
	_, err := os.Stat(pathStr)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("No such file.")
		} else if os.IsPermission(err) {
			return errors.New("Permission denied !")
		} else {
			return err
		}
	}
	return nil
}
