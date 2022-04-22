package api

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	"github.com/ganawaj/ipcalc/internal/iputil"
)

// Request format
type Request struct {
	Address string `json:"address"`
	Netmask string `json:"mask"`
}

// Response format
type Response struct {
	Address			string `json:"address"`
	Netmask			string `json:"mask"`
	Network 		string `json:"network"`
	Broadcast		string `json:"broadcast,omitempty"`
	HostMin			string `json:"hostmin,omitempty"`
	HostMax			string `json:"hostmax,omitempty"`
	Error 			string `json:"error,omitempty"`
}

// Server network configuration
type ServerConfig struct {
	Port    int
	Address string
}

// Server configuration
type Server struct {
	Config	ServerConfig
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
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


// create Response from request header
func (s *Server) newResponse(r *http.Request) (Response, error){

	request := Request{}

	// handle json request errors
	err := s.readJSON(r, &request)
	if err != nil {
		return Response{}, err
	}

	// get broadcast address
	broadcast_addr, err := iputil.GetBroadcastAddr(request.Address, request.Netmask)
	if err != nil {
		return Response{}, err
	}

	// get minimum host ip
	host_min, err := iputil.GetHostMin(request.Address, request.Netmask)
	if err != nil {
		return Response{}, err
	}

	// get maximum host ip
	host_max, err := iputil.GetHostMax(request.Address, request.Netmask)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Address: request.Address,
		Netmask: request.Netmask,
		Network: fmt.Sprintf("%s/%s", request.Address, request.Netmask),
		Broadcast: broadcast_addr,
		HostMin: host_min,
		HostMax: host_max,
	}, nil

}

// handles requests to '/'
func (s *Server) ipcalcHandler(w http.ResponseWriter, r *http.Request) {

	response, err := s.newResponse(r)
	if err != nil {
		s.ServerError(w, err)
		return
	}

	// handle json marshalling and marshalling errors
	b, err := json.Marshal(response)
	if err != nil {
		s.ServerError(w, err)
		return
	}

	b = append(b, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}