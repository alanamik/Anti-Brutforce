package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alanamik/Anti-Brutforce/internal/antibrutforce"
)

// A handler for test middleware.
func (s *Server) HelloHandler(w http.ResponseWriter, _ *http.Request) {
	res := ResponseSuccess{
		StatusCode: 200,
		Status:     "Access is allowed",
	}
	respondWithSuccess(w, res)
}

func (s *Server) AddWhiteIPHandler(w http.ResponseWriter, r *http.Request) {
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
	respondWithSuccess(w, res)
}

func (s *Server) AddBlackIPHandler(w http.ResponseWriter, r *http.Request) {
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
	err = s.Abf.AddToList(req.Cidr, false)
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
	respondWithSuccess(w, res)
}

func (s *Server) DeleteWhiteIPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondWithError(w, 500, "The method should be DELETE")
		return
	}
	var req DeleteIPIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 500, ErrInternalServerError500.Error())
		return
	}
	err := s.Abf.DeleteFromList(req.Cidr)
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
	respondWithSuccess(w, res)
}

func (s *Server) DeleteBlackIPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondWithError(w, 500, "The method should be DELETE")
		return
	}
	var req DeleteIPIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 500, ErrInternalServerError500.Error())
		return
	}
	err := s.Abf.DeleteFromList(req.Cidr)
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
	respondWithSuccess(w, res)
}

func (s *Server) ClearBucket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondWithError(w, 500, "The method should be DELETE")
		return
	}
	var req ClearBucketIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 500, ErrInternalServerError500.Error())
		return
	}
	err := s.Abf.ClearIPBuckets(req.IP)
	if errors.Is(err, antibrutforce.ErrNoSuchIP) {
		respondWithError(w, 400, antibrutforce.ErrNoSuchIP.Error())
		return
	}

	err = s.Abf.ClearLoginBuckets(req.Login)
	if errors.Is(err, antibrutforce.ErrNoSuchLogin) {
		respondWithError(w, 400, antibrutforce.ErrNoSuchLogin.Error())
		return
	}
	res := ResponseSuccess{
		StatusCode: 200,
		Status:     SuccessClearBucket,
	}
	respondWithSuccess(w, res)
}
