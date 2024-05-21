package handlers

import (
	"OTUS_hws/Anti-BruteForce/internal/gen/models"
	"OTUS_hws/Anti-BruteForce/internal/gen/restapi/operations"
	"OTUS_hws/Anti-BruteForce/internal/redisdb"

	"github.com/go-openapi/runtime/middleware"
)

func (h *Handler) AddToBL(params operations.BlacklistPutParams) middleware.Responder {

	err := h.abf.RedisServer.AddToList(params.HTTPRequest.Context(), *params.Body.IP, redisdb.Blacklist)
	if err != nil {
		return operations.NewBlacklistPutInternalServerError().WithPayload(&models.Error500{
			Error: err.Error(),
		},
		)
	}
	return operations.NewBlacklistPutOK().WithPayload(&models.Status200{Status: SuccessAddedStatus})
}
