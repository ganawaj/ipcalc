package api

import (
	"encoding/json"
	"net/http"
)


// handles requests to '/'
func (s *Server) ipcalcHandler(w http.ResponseWriter, r *http.Request) {

	request, err := s.newRequest(r)
	if err != nil {
		s.ServerError(w, err)
		return
	}

	response, err := s.newResponse(request)
	if err != nil {
		s.ServerError(w, err)
		return
	}

	// handle json marshalling and marshalling errors
	b, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		s.ServerError(w, err)
		return
	}

	b = append(b, '\n')

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(b)

}