package server

import (
	"apiscrape/lib/loggr"
	"apiscrape/omgr/data"
	"net/http"
)

type ApiServer struct {
	Server *http.Server
	Logger *loggr.Logger
}

func NewApiServer(store *data.Store, logger *loggr.Logger) *ApiServer {
	as := new(ApiServer)
	as.Server = new(http.Server)
	as.Logger = logger
	as.Server.Handler = NewHandler(store, logger)
	return as
}

func (s *ApiServer) Start(addr string) error {
	s.Server.Addr = addr
	s.Logger.Log("starting server on", addr)
	err := s.Server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *ApiServer) Stop() error {
	s.Logger.Log("stopping server")
	err := s.Server.Close()
	if err != nil {
		return err
	}

	return nil
}
