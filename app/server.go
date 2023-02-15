package app

import (
	"fmt"
	"hash/maphash"
	"net"
	"net/http"
	"time"

	"github.com/RobleDev498/spaces/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

func (c *CORSEnabler) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Expose-Headers", "Token")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Requested-With, X-Csrf-Token")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		h.ServeHTTP(w, r)
	})
}

func (s *Server) Start() error {
	//addr := "127.0.0.1:8080"
	host := s.Config.Server.Host
	port := s.Config.Server.Port
	addr := fmt.Sprintf("%s:%s", host, port)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowCredentials: true,
		ExposedHeaders:   []string{"Token"},
		AllowedHeaders:   []string{"Authorization", "X-Requested-With", "X-Csrf-Token"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	// c := NewCORSEnabler()
	handler := c.Handler(s.RootRouter)

	s.Server = &http.Server{
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         addr,
	}

	s.Server.ListenAndServe()

	return nil
}

func (a *App) OriginChecker() func(*http.Request) bool {
	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}

		if origin == "http://localhost:3000" {
			return true
		}

		return false
		//return true
	}
}

func (s *Server) getServer() *http.Server {
	return s.Server
}
