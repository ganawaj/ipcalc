package api

import (
	"net/http"
	"fmt"
	"runtime/debug"
	"encoding/json"
	"errors"
	"io"
)

func (s *Server) ServerError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	s.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// To be implemented.
func (s *Server) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {

			case errors.As(err, &syntaxError):
				return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

			case errors.Is(err, io.ErrUnexpectedEOF):
				return errors.New("body contains badly-formed JSON")

			case errors.As(err, &unmarshalTypeError):
				if unmarshalTypeError.Field != "" {
					return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
				}
				return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

			case errors.Is(err, io.EOF):
				return errors.New("body must not be empty")

			case errors.As(err, &invalidUnmarshalError):
					panic(err)

			default:
				return err
		}
	}

	return nil
}
