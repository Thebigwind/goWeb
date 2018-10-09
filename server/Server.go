package server

import (
	"net/http"
)

type RestServer struct {
	addr string
}

var GlobalRestServer *RestServer = nil

func NewRESTServer(addr string) *RestServer {
	GlobalRestServer = &RestServer{
		addr: addr,
	}
	return GlobalRestServer
}

func (server *RestServer) StartRESTServer() {
	router := NewRouter()
	//GetGlobeLocker()
	http.ListenAndServe(server.addr, router) //这个函数指这定地址和对应的handler
}
