package handler

import (
	"errors"
	"fmt"

	"net/http"
	"regexp"
	"strconv"

	. "github.com/bigwind/goWeb/common"
	. "github.com/bigwind/goWeb/server/manager"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("redirect")
	http.Redirect(w, r, "/main.html", http.StatusFound)
	return
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
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	ret := &ReturnJson{}
	defer r.Body.Close()
	/*
		Request本身也提供了FormValue()函数来获取用户提交的参数。如r.Form["username"]也可写成r.FormValue("username")。
		调用r.FormValue时会自动调用r.ParseForm，所以不必提前调用。
		r.FormValue只会返回同名参数中的第一个，若参数不存在则返回空字符串。
	*/

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
		fmt.Println("redirct before:")
		http.Redirect(w, r, "/main.html", http.StatusFound) //main.html
		fmt.Println("redirct :")
		return

	} else {
		err = errors.New("Empty username or password!")
		Logger.Errorf(err.Error())
		ret.Status = XT_API_RET_ERROR
		ret.Errmsg = err.Error()
		ret.AjaxReturnJson(w)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ret := &ReturnJson{}
	defer r.Body.Close()
	/*
		Request本身也提供了FormValue()函数来获取用户提交的参数。如r.Form["username"]也可写成r.FormValue("username")。
		调用r.FormValue时会自动调用r.ParseForm，所以不必提前调用。
		r.FormValue只会返回同名参数中的第一个，若参数不存在则返回空字符串。
	*/

	err := r.ParseForm()
	if err != nil {
		err = errors.New("Invalid form format!")
		Logger.Errorf(err.Error())
		ret.Status = XT_API_RET_ERROR
		ret.Errmsg = err.Error()
		ret.AjaxReturnJson(w)
		return
	}

	if len(r.Form["username"][0]) == 0 {
		//为空的处理
	}
	if len(r.Form["password"][0]) == 0 {
		//为空的处理
	}
	if len(r.FormValue("password")) < 6 {
		//密码长度处理
	}
	if r.FormValue("password") != r.FormValue("repassword") {
		//密码不一致处理
	}

	//年龄校验
	getint, err := strconv.Atoi(r.Form.Get("age"))
	if err != nil {
		//数字转化出错了，那么可能就不是数字
		return
	}
	//接下来就可以判断这个数字的大小范围了
	if getint > 100 {
		//太大了
	}
	if getint < 1 {
		//太小了
	}

	//中文校验
	if m, _ := regexp.MatchString("^\\p{Han}+$", r.Form.Get("realname")); !m {
		return
	}

	//英文校验
	if m, _ := regexp.MatchString("^[a-zA-Z]+$", r.Form.Get("engname")); !m {
		return
	}

	//性别按钮处理
	/*
		<input type="radio" name="gender" value="1">男
		<input type="radio" name="gender" value="2">女
	*/
	slice := []string{"1", "2"}

	for _, v := range slice {
		if v == r.Form.Get("gender") {
			return
		}
	}

	//复选框按钮处理
	/*
		<input type="checkbox" name="interest" value="football">足球
		<input type="checkbox" name="interest" value="basketball">篮球
		<input type="checkbox" name="interest" value="tennis">网球
	*/
	/*
		slicebox := []string{"football", "basketball", "tennis"}
		a := Slice_diff(r.Form["interest"], slicebox)
		if a == nil {
			return
		}
	*/

	//email 校验
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, r.Form.Get("email")); !m {
		fmt.Println("no")
	} else {
		fmt.Println("yes")
	}
	//手机号校验
	if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, r.Form.Get("mobile")); !m {
		return
	}

	//身份证号校验
	//验证15位身份证，15位的是全部数字
	if m, _ := regexp.MatchString(`^(\d{15})$`, r.Form.Get("usercard")); !m {
		return
	}

	//验证18位身份证，18位前17位为数字，最后一位是校验位，可能为数字或字符X。
	if m, _ := regexp.MatchString(`^(\d{17})([0-9]|X)$`, r.Form.Get("usercard")); !m {
		return
	}
}
