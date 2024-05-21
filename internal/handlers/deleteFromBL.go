package handlers

import (
	"OTUS_hws/Anti-BruteForce/internal/gen/models"
	"OTUS_hws/Anti-BruteForce/internal/gen/restapi/operations"
	"OTUS_hws/Anti-BruteForce/internal/redisdb"

	"github.com/go-openapi/runtime/middleware"
)

func (h *Handler) DeleteFromBL(params operations.BlacklistDeleteParams) middleware.Responder {

	err := h.abf.RedisServer.DeleteFromList(params.HTTPRequest.Context(), *params.Body.IP, redisdb.Blacklist)
	if err != nil {
		return operations.NewBlacklistDeleteInternalServerError().WithPayload(&models.Error500{
			Error: err.Error(),
		},
		)
	}

	return operations.NewBlacklistDeleteOK().WithPayload(&models.Status200{Status: SuccessDeletedStatus})
}
