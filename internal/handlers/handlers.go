package handlers

import (
	"OTUS_hws/Anti-BruteForce/internal/antibrutforce"
	"OTUS_hws/Anti-BruteForce/internal/gen/restapi/operations"
)

type Handler struct {
	abf *antibrutforce.AntiBrutForce
}

func NewHandler(abf *antibrutforce.AntiBrutForce) *Handler {
	return &Handler{
		abf: abf,
		//redisServer: redisS,
		//clients:     make(map[string]Bucket),
	}
}

func (h *Handler) Register(api *operations.AntiBrutForceAPI) {
	api.AuthCheckHandler = operations.AuthCheckHandlerFunc(h.AuthCheck)
	api.BlacklistDeleteHandler = operations.BlacklistDeleteHandlerFunc(h.DeleteFromBL)
	api.BlacklistPutHandler = operations.BlacklistPutHandlerFunc(h.AddToBL)
}
