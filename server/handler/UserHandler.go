package handler

import (
	"errors"
	"fmt"

	"net/http"

	. "github.com/bigwind/goWeb/common"
	. "github.com/bigwind/goWeb/server/manager"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("redirect")
	http.Redirect(w, r, "/main.html", http.StatusFound)
}

/*
func RootHandler2(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		fmt.Println("userName:", userName)
		switch userName {
		case "":
			log.Fatal("user already logout!")
			http.Redirect(w, r, "/login.html", http.StatusFound)
		default:
			http.Redirect(w, r, "/dashboard.html", http.StatusFound)
		}
	} else {
		//If userName is not a string typeauthenticated-user-session
		log.Fatal("cookie on client not found!")
		http.Redirect(w, r, "/login.html", http.StatusFound)
	}
}
*/
func Login(w http.ResponseWriter, r *http.Request) {
	ret := &ReturnJson{}
	err := r.ParseForm()
	if err != nil {
		err = errors.New("Invalid form format!")
		Logger.Errorf(err.Error())
		ret.Status = XT_API_RET_ERROR
		ret.Errmsg = err.Error()
		ret.AjaxReturnJson(w)
		return
	}

	if r.FormValue("username") != "" && r.FormValue("password") != "" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		mgr := GetUserMgr()
		err := mgr.CheckLoginUserPasswd(username, password, "frontend")
		if err != nil {
			err = errors.New("User not exist or inactive!")
			Logger.Errorf(err.Error())
			ret.Status = XT_API_RET_ERROR
			ret.Errmsg = err.Error()
			ret.AjaxReturnJson(w)
			return
		}
		//用户信息存入session
		session, _ := sessionStore.New(r, "authenticated-user-session")
		session.Options.MaxAge = 86400 * 15
		session.Values["username"] = username

		mgr = GetUserMgr()
		err, userID := mgr.GetUserId(username)
		if err != nil {
			err = errors.New("User info error!")
			Logger.Errorf(err.Error())
			ret.Status = XT_API_RET_ERROR
			ret.Errmsg = err.Error()
			ret.AjaxReturnJson(w)
			return
		}

		Logger.Infoln("GetUserID id=%d", userID)
		session.Values["userid"] = userID
		err = session.Save(r, w)
		if err != nil {
			err = errors.New("Set user session/cookie error!")
			Logger.Errorf(err.Error())
			ret.Status = XT_API_RET_ERROR
			ret.Errmsg = err.Error()
			ret.AjaxReturnJson(w)
			return
		}

		//页面跳转
		//Logger.Infoln("auth pass and session saved!", session)
		//http.Redirect(w, r, "/index.html", http.StatusFound)
		http.Redirect(w, r, "/main.html", http.StatusFound)

	} else {
		err = errors.New("Empty username or password!")
		Logger.Errorf(err.Error())
		ret.Status = XT_API_RET_ERROR
		ret.Errmsg = err.Error()
		ret.AjaxReturnJson(w)
		return
	}
}
