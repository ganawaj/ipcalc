package api

import(
	"fmt"
	"net/http"

	"github.com/ganawaj/ipcalc/internal/iputil"
)

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

// create Response from request header
func (s *Server) newResponse(r *http.Request) (Response, error){

	request := Request{}

	// handle json request errors
	err := s.readJSON(r, &request)
	if err != nil {
		return Response{}, err
	}

	uip := iputil.New(request.Address, request.Netmask)

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
		Address: request.Address,
		Netmask: request.Netmask,
		Network: fmt.Sprintf("%s/%s", request.Address, request.Netmask),
		Broadcast: broadcast_addr,
		HostMin: host_min,
		HostMax: host_max,
	}, nil

}
