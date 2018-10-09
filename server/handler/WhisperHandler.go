package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	. "github.com/xtao/goWeb/common"
	. "github.com/xtao/goWeb/server/manager"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("redirect")
	http.Redirect(w, r, "/main.html", http.StatusFound)
}

func RunCmdHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	err := req.ParseForm()
	ret := &ReturnJson{
		Status: XT_API_RET_OK,
	}
	if err != nil {
		ret.Status = XT_API_RET_ERROR
		ret.Errmsg = err.Error()
		ret.AjaxReturnJson(w)
		return
	}

	startTime := strings.TrimSpace(req.Form.Get("startTime"))
	if startTime == "" {
		startTime = "2000-01-01 00:00:00"
		fmt.Println("startTime:", startTime)
	}

	endTime := strings.TrimSpace(req.Form.Get("endTime"))
	if endTime == "" {
		endTime = "2100-01-01 00:00:00"
		fmt.Println("endTime:", endTime)
	}

	duration := strings.TrimSpace(req.Form.Get("duration"))
	tm, err := strconv.Atoi(duration)
	if err != nil {
		ret.Status = XT_API_RET_ERROR
		ret.Errmsg = "Invalid value."
		ret.AjaxReturnJson(w)
		return
	}
	if tm%10 != 0 {
		tm = (tm / 10) * 10
	}
	if tm < 10 {
		tm = 10
	}
	fmt.Println(PROG_PATH)
	mgr := &ExecMgr{
		Prog: "python", //PROG_PATH
		//		Args: []string{strconv.Itoa(tm)},
		Args: []string{"../whisper_export.py", strconv.Itoa(tm), startTime, endTime}, //"http://www.xtaotech.com/image/icon.png", "-o", "../Upload/icon.png"
	}
	fmt.Println("mgr:", mgr)
	err = mgr.ExecCmd()
	if err != nil {
		ret.Status = XT_API_RET_ERROR
		ret.Errmsg = err.Error()
		ret.AjaxReturnJson(w)
		return
	}
	ret.Result = tm
	ret.AjaxReturnJson(w)
	return
}

func DownLoadFileHandler(w http.ResponseWriter, req *http.Request) {
	//locker := GetGlobeLocker()
	//locker.RLock()
	//defer locker.RUnlock()
	defer req.Body.Close()

	vars := mux.Vars(req)
	duration := vars["duration"]
	startTime := vars["startTime"]
	endTime := vars["endTime"]

	if startTime == "" {
		startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05")
		fmt.Println("startTime:", startTime)
	}
	if endTime == "" {
		endTime = time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("endTime:", endTime)
	}
	if duration == "" {
		duration = "30"
	}

	tm, err := strconv.Atoi(duration)
	if err != nil {
		err = errors.New("the duration format not correct. it should be a number, like 10, 20, 30")
		w.Header().Set("Content-Type", " text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		err = errors.New("the startTime format not correct. it should like '2006-01-02 15:04:05'")
		w.Header().Set("Content-Type", " text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		err = errors.New("the endTime format not correct. it should like '2006-01-02 15:04:05' ")
		w.Header().Set("Content-Type", " text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if tm%10 != 0 {
		tm = (tm / 10) * 10
	}
	if tm < 10 {
		tm = 10
	}
	//fmt.Println(PROG_PATH)
	mgr := &ExecMgr{
		Prog: "python", //PROG_PATH
		//		Args: []string{strconv.Itoa(tm)},
		Args: []string{"../whisper_export.py", strconv.Itoa(tm), startTime, endTime}, //"http://www.xtaotech.com/image/icon.png", "-o", "../Upload/icon.png"
	}
	fmt.Println("mgr:", mgr)
	err = mgr.ExecCmd()

	if err != nil {
		w.Header().Set("Content-Type", " text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//ret.Result = tm
	//ret.AjaxReturnJson(w)
	//return
	fpath := "/var/whisper_data_result.txt"
	err = CheckReadFile(fpath)
	if err != nil {
		w.Header().Set("Content-Type", " text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	file, err := os.Open(fpath)
	if err != nil {
		w.Header().Set("Content-Type", " text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	fileName := path.Base(fpath)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	fstat, _ := file.Stat()
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fstat.Size()))
	io.Copy(w, file)
}
