package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// Server アプリ本体機能を持つ
type Server struct {
	srv *http.Server
}

// New Server を生成
func New(port int) *Server {
	srv := &Server{
		srv: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
			// Good practice to set timeouts to avoid Slowloris attacks.（ここの値は検証用なので注意）
			// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
			WriteTimeout: time.Second * 25,
			ReadTimeout:  time.Second * 25,
			IdleTimeout:  time.Second * 60,
		},
	}
	srv.routes()
	return srv
}

// Run Server 起動
func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

// ServeHTTP これを定義することで、Server 自体が http.Handler interface を満たすことができる
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.srv.Handler.ServeHTTP(w, r)
}

func (s *Server) routes() {
	router := mux.NewRouter()
	router.Handle("/health", s.healthy()).Methods(http.MethodGet)
	router.Handle("/metrics", s.healthy()).Methods(http.MethodGet)
	s.srv.Handler = router
}

// Shutdown shutdown を停止
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// healthy /health
func (s *Server) healthy() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`"ok"`))
	})
}

func main() {
	srv := New(8888)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server is ready at", ":8888")
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server on :%d: %v", 8888, err)
		}
	}()

	<-quit

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("could not gracefully shutdown the server: %v", err)
	}
	log.Println("Server stopped")
}
