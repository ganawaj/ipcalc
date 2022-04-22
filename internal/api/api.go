package api

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	"github.com/ganawaj/ipcalc/internal/iputil"
)


type Request struct {
	Address string `json:"address"`
	Netmask string `json:"mask"`
}

type Response struct {
	Address			string `json:"address"`
	Netmask			string `json:"mask"`
	Network 		string `json:"network"`
	Broadcast		string `json:"broadcast,omitempty"`
	HostMin			string `json:"hostmin,omitempty"`
	HostMax			string `json:"hostmax,omitempty"`
	Error 			string `json:"error,omitempty"`
}

type ServerConfig struct {
	Port    int
	Address string
}

type Server struct {
	Config	ServerConfig
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
}

func (s *Server) ListenAndServe() error {

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%d", s.Config.Address, s.Config.Port),
		ErrorLog: s.ErrorLog,
		Handler:  s.Routes(),
	}

	s.InfoLog.Printf("Starting server on http://%s:%d", s.Config.Address, s.Config.Port)

	return server.ListenAndServe()

}

func (s *Server) Routes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ipcalcHandler)

	return s.recoverPanic(s.logRequest(mux))
}

func (s *Server) newResponse(r *http.Request) (Response, error){



	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return Response{}, err
	}

	broadcast_addr, err := iputil.GetBroadcastAddr(request.Address, request.Netmask)
	if err != nil {
		return Response{}, err
	}

	host_min, err := iputil.GetHostMin(request.Address, request.Netmask)
	if err != nil {
		return Response{}, err
	}

	host_max, err := iputil.GetHostMax(request.Address, request.Netmask)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		Address: request.Address,
		Netmask: request.Netmask,
		Network: fmt.Sprintf("%s/%s", request.Address, request.Netmask),
		Broadcast: broadcast_addr,
		HostMin: host_min,
		HostMax: host_max,
	}

	return response, nil

}

func (s *Server) ipcalcHandler(w http.ResponseWriter, r *http.Request) {

	request := Request{}

	err := s.readJSON(w, r, request)
	if err != nil {
		s.ServerError(w, err)
		return
	}
	

	b, err := json.Marshal(response)
	if err != nil {
		s.ServerError(w, err)
		return
	}

	b = append(b, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}