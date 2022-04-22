package api

import(
	"log"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Server configuration
type Server struct {
	Config	ServerConfig
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
}

// Server network configuration
type ServerConfig struct {
	Port    int
	Address string
}

// Setup server and listen on address and ports
func (s *Server) ListenAndServe() error {

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%d", s.Config.Address, s.Config.Port),
		ErrorLog: s.ErrorLog,
		Handler:  s.Routes(),
	}

	s.InfoLog.Printf("Starting server on http://%s:%d", s.Config.Address, s.Config.Port)

	return server.ListenAndServe()

}

// Sets and creates routes and returns handler
func (s *Server) Routes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ipcalcHandler)

	// handle errors and recover requests
	return s.recoverPanic(s.logRequest(mux))
}

func (s *Server) ServerError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	s.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
