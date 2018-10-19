package server

import (
	"net/http"

	. "github.com/bigwind/goWeb/common"
	. "github.com/bigwind/goWeb/server/handler"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func NewRouter() *negroni.Negroni {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", RootHandler).Methods("GET")
	router.HandleFunc("/v1/goweb/login", Login).Methods("POST")
	//router.HandleFunc("/api/cmd/run", RunCmdHandler).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(STATIC_DIR)))
	n := negroni.Classic()
	n.UseHandler(router)
	return n
}
