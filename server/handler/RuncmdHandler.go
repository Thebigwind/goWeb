package handler

import (
	//"errors"
	"fmt"
	//"io"
	"net/http"

	//"os"
	//"path"
	"strconv"
	"strings"

	//"time"

	//"github.com/gorilla/mux"
	. "github.com/bigwind/goWeb/common"
	. "github.com/bigwind/goWeb/server/manager"
)

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

	fmt.Println(PROG_PATH)
	mgr := &ExecMgr{
		Prog: "python",                                                               //PROG_PATH
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
