package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/sessions"
)

//var sessionStore *PGStore

var sessionStore, _ = pgstore.NewPGStore("postgres://postgres:123456@127.0.0.1:5432/testdb?sslmode=disable", []byte("secret-key")) //verify-full

func GetSession(r *http.Request) (*sessions.Session, error) {
	return sessionStore.Get(r, "authenticated-user-session")
}

func GetUserId32(req *http.Request) int32 {
	sess, err := GetSession(req)
	if err != nil {
		return -1
	}
	uid, ok := sess.Values["userid"]
	if !ok {
		return -1
	}
	return int32(uid.(int))
}

func GetUserIdInt(req *http.Request) int {
	sess, err := GetSession(req)
	if err != nil {
		return -1
	}
	uid, ok := sess.Values["userid"]
	if !ok {
		return -1
	}
	return uid.(int)
}

func GetUserName(req *http.Request) string {
	sess, err := GetSession(req)
	if err != nil {
		return ""
	}
	uname, ok := sess.Values["username"]
	if !ok {
		return ""
	}
	return uname.(string)
}

func init() {
	fmt.Println("xxxxxxx")
	sessionStore.Options = &sessions.Options{
		Domain:   "127.0.0.1",
		Path:     "/",
		MaxAge:   86400 * 15, //seconds
		HttpOnly: false,
	}

	//defer sessionStore.Close()
	sessionStore.Cleanup(time.Minute * 15)
}
