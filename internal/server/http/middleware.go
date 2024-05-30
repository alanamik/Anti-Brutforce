package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) CheckRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range HandlersPaths {
			if r.URL.Path == h {
				next.ServeHTTP(w, r)
				return
			}
		}
		var req CheckRequestIn
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, 500, ErrInternalServerError500.Error())
			return
		}
		pass, err := s.Abf.CheckRequest(r.Context(), req.IP, req.Login, req.Password)
		if err != nil {
			respondWithError(w, 500, ErrInternalServerError500.Error())
			return
		}
		if !pass {
			respondWithError(w, 400, ErrAccessDenied)
			return
		}
		next.ServeHTTP(w, r)
	})
}