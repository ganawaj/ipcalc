package api

import (
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






func (s *Server) newRequest(r *http.Request) (Request, error){

	request := Request{}

	// handle json request errors
	err := s.readJSON(r, &request)
	if err != nil {
		return Request{}, err
	}

	return request, nil
}


// create Response from request header
func (s *Server) newResponse(r Request) (Response, error){

	uip := iputil.New(r.Address, r.Netmask)

	// get broadcast address
	broadcast_addr, err := uip.GetBroadcastAddr()
	if err != nil {
		return Response{}, err
	}

	// get minimum host ip
	host_min, err := uip.GetHostMin()
	if err != nil {
		return Response{}, err
	}

	// get maximum host ip
	host_max, err := uip.GetHostMax()
	if err != nil {
		return Response{}, err
	}

	return Response{
		Address: r.Address,
		Netmask: r.Netmask,
		Network: fmt.Sprintf("%s/%s", r.Address, r.Netmask),
		Broadcast: broadcast_addr,
		HostMin: host_min,
		HostMax: host_max,
	}, nil

}
