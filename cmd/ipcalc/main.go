package ipcalc

import (
	"log"
	"net/http"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"inet.af/netaddr"
	"github.com/go-kit/kit/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
)

var ErrEmpty = errors.New("Empty string")

type Service interface {
	GetNetworkInfo(string, string) (IP, error)
}

type service struct {}

type IP struct {
	Address			string `json:"address"`
	Netmask			string `json:"mask"`
	Network 		string `json:"network"`
	Broadcast		string `json:"broadcast"`
	HostMin			string `json:"hostmin"`
	HostMax			string `json:"hostmax"`
	Error 			string `json:"error,omitempty"`
}

type NetworkRequest struct {
	Address string `json:"address"`
	Netmask string `json:"mask"`
}

func makeNetworkEndpoint(svc Service) endpoint.Endpoint {

	return func(_ context.Context, request interface{}) (interface{}, error) {

		req := request.(NetworkRequest)

		net, err := svc.GetNetworkInfo(req.Address, req.Netmask)
		if err != nil {
			return IP{}, nil
		}
		fmt.Println(net)

		return net, nil
	}
}

func decodeNetworkAddressRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request NetworkRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main(){

	svc := service{}

	networkaddressHandler := httptransport.NewServer(
		makeNetworkEndpoint(svc),
		decodeNetworkAddressRequest,
		encodeResponse,

	)

	http.Handle("/network", networkaddressHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (service) GetNetworkInfo (addr string, mask string) (IP, error){

	if addr == "" || mask == "" {
		return IP{}, ErrEmpty
	}

	// create prefix from ip address and range
	prefix, err := netaddr.ParseIPPrefix(fmt.Sprintf("%s/%s", addr, mask))
	if err != nil {
		return IP{}, ErrEmpty
	}

	return IP{
		Address: 	addr,
		Netmask: 	mask,
		Network: 	fmt.Sprintf("%s/%s", prefix.Range().From().String(), mask),
		Broadcast: 	prefix.Range().To().String(),
		HostMin: 	prefix.Range().From().Next().String(),
		HostMax: 	prefix.Range().To().Prior().String(),
	}, nil

}