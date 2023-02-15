package app

import (
	"fmt"
	"hash/maphash"
	"net"
	"net/http"
	"time"

	"github.com/RobleDev498/spaces/config"
	"github.com/gorilla/mux"
)

type Server struct {
	RootRouter *mux.Router
	Server     *http.Server
	ListenAddr *net.TCPAddr

	LocalRouter *mux.Router
	Config      *config.Config

	hub *Hub

	hashSeed maphash.Seed
}

func NewServer() (*Server, error) {
	rootRouter := mux.NewRouter()

	s := &Server{
		RootRouter: rootRouter,
		hashSeed:   maphash.MakeSeed(),
	}

	//s.LocalRouter = s.RootRouter.PathPrefix("/static").Subrouter()
	fs := http.FileServer(http.Dir("./static"))
	s.RootRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	return s, nil
}

func (a *App) SetHub(hub *Hub) {
	a.srv.hub = hub
}

type CORSEnabler struct{}

func NewCORSEnabler() CORSEnabler {
	return CORSEnabler{}
}

func (s *Server) Start() error {
	//addr := "127.0.0.1:8080"
	host := s.Config.Server.Host
	port := s.Config.Server.Port
	addr := fmt.Sprintf("%s:%s", host, port)

	s.Server = &http.Server{
		Handler:      s.RootRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         addr,
	}

	s.Server.ListenAndServe()

	return nil
}

func (s *Server) getServer() *http.Server {
	return s.Server
}
