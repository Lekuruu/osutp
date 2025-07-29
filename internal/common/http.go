package common

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server is the main struct that holds the state for an http server.
type Server struct {
	Host      string
	Port      int
	Name      string
	State     *State
	Router    *mux.Router
	StaticDir string
}

func NewServer(host string, port int, name string, state *State) *Server {
	return &Server{
		Host:   host,
		Port:   port,
		Name:   name,
		State:  state,
		Router: mux.NewRouter(),
	}
}

// Context is a struct that holds the request context for each endpoint call.
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	State    *State
	Vars     map[string]string
}

// Serve starts the HTTP server and listens for incoming requests.
func (server *Server) Serve() {
	bind := fmt.Sprintf(
		"%s:%d",
		server.Host,
		server.Port,
	)
	log.Printf("Starting server on %s\n", bind)

	err := http.ListenAndServe(bind, server.LoggingMiddleware(server.Router))
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
		return
	}
}

// ResponseContext is a wrapper around http.ResponseWriter that
// allows us to capture the status code of a response.
type ResponseContext struct {
	w http.ResponseWriter
	s int
}

func (rc *ResponseContext) Header() http.Header {
	return rc.w.Header()
}

func (rc *ResponseContext) Write(b []byte) (int, error) {
	return rc.w.Write(b)
}

func (rc *ResponseContext) WriteHeader(status int) {
	rc.s = status
	rc.w.WriteHeader(status)
}

func (rc *ResponseContext) Status() int {
	if rc.s == 0 {
		return http.StatusOK
	}
	return rc.s
}

// ContextMiddleware creates a new Context struct for each request.
func (server *Server) ContextMiddleware(handler func(*Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := &Context{
			Response: w,
			Request:  r,
			State:    server.State,
			Vars:     mux.Vars(r),
		}

		w.Header().Set("Server", server.Name)
		handler(context)
	}
}

// LoggingMiddleware logs the details of each request.
func (server *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := &ResponseContext{w: w}
		start := time.Now()
		next.ServeHTTP(rc, r)
		duration := time.Since(start)
		log.Printf("[%d] %s %s (%s)",
			rc.Status(),
			r.Method,
			r.RequestURI,
			duration,
		)
	})
}
