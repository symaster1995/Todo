package http

import (
	"Todo"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net"
	"net/http"
	"strconv"
	"time"
)

const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln     net.Listener
	server *http.Server
	router chi.Router

	Addr   int

	TaskService Todo.TaskService
}

func (s *Server) Open() (err error) {

	if s.ln, err = net.Listen("tcp", ":" + strconv.Itoa(s.Addr)); err != nil {
		return err
	}

	go s.server.Serve(s.ln)

	fmt.Println("SERVED", s.Addr)

	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: chi.NewRouter(),
	}

	s.server.Handler = http.HandlerFunc(s.serveHttp)

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)

	s.router.Route("/v1", func(r chi.Router) {
		r.Use(apiVersionCtx("v1"))
		r.Mount("/task", s.mountTodoRoutes())
	})

	return s
}

func (s *Server) serveHttp(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "api.version", version))
			next.ServeHTTP(w, r)
		})
	}
}
