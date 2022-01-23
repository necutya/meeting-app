package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"server/config"
	"server/internal/server/handlers"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

const version1 = "/v1"

type Server struct {
	http   *http.Server
	config *config.HTTP

	// graceful shutdown stuff
	runErr    error
	readiness bool

	// handlers
	meetingHandler *handlers.MeetingHandler
}

func New(cfg *config.HTTP,
	meetingHandler *handlers.MeetingHandler,
) *Server {
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	}

	srv := Server{
		config:         cfg,
		meetingHandler: meetingHandler,
	}

	srv.setupHTTP(&httpSrv)
	return &srv
}

func (s *Server) setupHTTP(srv *http.Server) {
	handler := s.setupHandler()

	srv.Handler = handler
	s.http = srv
}
func (s *Server) setupHandler() http.Handler {
	var (
		router    = mux.NewRouter()
		apiRouter = router.PathPrefix(s.config.URLPrefix).Subrouter()
		v1Router  = apiRouter.PathPrefix(version1).Subrouter()

		publicChain = alice.New()
	)
	log.Info()
	// public routes
	v1Router.Handle("/create", publicChain.ThenFunc(s.meetingHandler.CreateRoom)).Methods(http.MethodGet)
	v1Router.Handle("/join/{roomID}", publicChain.ThenFunc(s.meetingHandler.JoinRoom)).Methods(http.MethodGet)

	return cors.New(cors.Options{
		AllowedOrigins: s.config.CORSAllowedHost,
		// AllowedMethods: []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut,
		// 	http.MethodDelete, http.MethodOptions, http.MethodPatch},
		// AllowedHeaders:     []string{"*"},
		// AllowCredentials:   true,
		// OptionsPassthrough: false,
	}).Handler(router)
	// return router
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Debug(fmt.Sprintf("HTTP server run (addr: %s)", s.http.Addr))
		err := s.http.ListenAndServe()
		s.runErr = err
		log.Info(fmt.Sprintf("HTTP server stop (err: %s)", err))
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)

		if err != nil {
			log.Info(fmt.Sprintf("HTTP server shutdown error (err: %s)", err))
		}
	}()

	s.readiness = true
}
