package server

import (
	"apiscrape/lib/loggr"
	"apiscrape/oapi/data"
	"net/http"
)

type ApiServer struct {
	Server *http.Server
	Logger *loggr.Logger
}

func NewApiServer(store *data.OApiStore, logger *loggr.Logger) *ApiServer {
	s := new(ApiServer)
	s.Server = new(http.Server)
	s.Logger = logger
	s.Server.Handler = NewHandler(store, logger)
	return s
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
