package server

import (
	"net/http"
	/*
	   golang自带的http.SeverMux路由实现简单,本质是一个map[string]Handler,是请求路径与该路径对应的处理函数的映射关系。实现简单功能也比较单一：

	   不支持正则路由， 这个是比较致命的
	   只支持路径匹配，不支持按照Method，header，host等信息匹配，所以也就没法实现RESTful架构
	   而gorilla/mux是一个强大的路由，小巧但是稳定高效，不仅可以支持正则路由还可以按照Method，header，host等信息匹配，
	   可以从我们设定的路由表达式中提取出参数方便上层应用，而且完全兼容http.ServerMux
	*/
	"github.com/gorilla/mux"
	. "github.com/xtao/webGo/common"
	. "github.com/xtao/webGo/server/handler"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router { //

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(STATIC_DIR)))

	return router
}

var routes = Routes{

	Route{
		"DownLoadFileHandler",
		"GET",
		"/api/file/download/{duration}/{startTime}/{endTime}",
		DownLoadFileHandler,
	},
}
