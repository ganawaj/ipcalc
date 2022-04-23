package api

import(
	"net/http"
	"encoding/json"
)

// handles requests to '/'
func (s *Server) ipcalcHandler(w http.ResponseWriter, r *http.Request) {

	response, err := s.newResponse(r)
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