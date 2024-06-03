package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrInternalServerError500 = errors.New("internal server error")
	SuccessAddListStatus      = "IP successfully added in the list"
	SuccessDeleteListStatus   = "IP successfully deleted from the list"
	SuccessClearBucket        = "Successfully deleted the bucket"
	ErrAccessDenied           = "Access denied"
)

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hello", s.CheckRequest(http.HandlerFunc(s.HelloHandler)))
	// mux.HandleFunc("/addWhiteIp", s.AddWhiteIPHandler)
	// mux.HandleFunc("/deleteWhiteIP", s.DeleteWhiteIPHandler)
	// mux.HandleFunc("/addBlackIp", s.AddBlackIPHandler)
	// mux.HandleFunc("/deleteBlackIP", s.DeleteBlackIPHandler)
	// mux.HandleFunc("/clearBucket", s.ClearBucket)

	mux.Handle("/addWhiteIp", s.CheckRequest(http.HandlerFunc(s.AddWhiteIPHandler)))
	mux.Handle("/deleteWhiteIP", s.CheckRequest(http.HandlerFunc(s.DeleteWhiteIPHandler)))
	mux.Handle("/addBlackIp", s.CheckRequest(http.HandlerFunc(s.AddBlackIPHandler)))
	mux.Handle("/deleteBlackIP", s.CheckRequest(http.HandlerFunc(s.DeleteBlackIPHandler)))
	mux.Handle("/clearBucket", s.CheckRequest(http.HandlerFunc(s.ClearBucket)))

	return mux
}

func respondWithSuccess(w http.ResponseWriter, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	res := ResponseError{
		StatusCode: code,
		Error:      msg,
	}
	response, err := json.Marshal(res)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}
