package server

import (
	"OTUS_hws/Anti-BruteForce/internal/antibrutforce"
	"encoding/json"
	"errors"
	"fmt"
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
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/HELLO")
	})

	mux.HandleFunc("/addWhiteIp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, 500, "The method should be POST")
			return
		}
		var req AddIPIn
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			respondWithError(w, 500, ErrInternalServerError500.Error())
			return
		}
		err = s.Abf.AddToList(req.Cidr, true)
		if err != nil {
			if errors.Is(err, antibrutforce.ErrIPInListYet) {
				respondWithError(w, 400, err.Error())
			} else {
				respondWithError(w, 500, ErrInternalServerError500.Error())
			}
			return
		}

		res := ResponseSuccess{
			StatusCode: 200,
			Status:     SuccessAddListStatus,
		}
		respondWithSuccess(w, 200, res)
	})

	mux.HandleFunc("/deleteWhiteIP", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respondWithError(w, 500, "The method should be DELETE")
			return
		}
		var req DeleteIPIn
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, 500, ErrInternalServerError500.Error())
			return
		}
		err := s.Abf.DeleteFromList(req.IP)
		if err != nil {
			if errors.Is(err, antibrutforce.ErrNoSuchIP) {
				respondWithError(w, 400, antibrutforce.ErrNoSuchIP.Error())
			} else {
				respondWithError(w, 500, ErrInternalServerError500.Error())
			}
			return
		}

		res := ResponseSuccess{
			StatusCode: 200,
			Status:     SuccessDeleteListStatus,
		}
		respondWithSuccess(w, 200, res)
	})

	mux.HandleFunc("/clearBucket", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respondWithError(w, 500, "The method should be DELETE")
			return
		}
		var req ClearBucketIn
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, 500, ErrInternalServerError500.Error())
			return
		}
		// err := s.Abf.ClearIPBuckets(req.Ip)
		// if errors.Is(err, antibrutforce.ErrNoSuchIP) {
		// 	respondWithError(w, 400, antibrutforce.ErrNoSuchIP.Error())
		// 	return
		// }
		err := s.Abf.ClearLoginBuckets(req.Login)
		if errors.Is(err, antibrutforce.ErrNoSuchLogin) {
			respondWithError(w, 400, antibrutforce.ErrNoSuchLogin.Error())
			return
		}

		res := ResponseSuccess{
			StatusCode: 200,
			Status:     SuccessClearBucket,
		}
		respondWithSuccess(w, 200, res)
	})

	return mux
}

func respondWithSuccess(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
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
