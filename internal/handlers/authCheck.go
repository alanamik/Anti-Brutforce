package handlers

import (
	"OTUS_hws/Anti-BruteForce/internal/gen/restapi/operations"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
)

func (h *Handler) AuthCheck(params operations.AuthCheckParams) middleware.Responder {
	fmt.Println(*params.Body.IP, *params.Body.Login, *params.Body.Password)
	pass, err := h.abf.CheckRequest(params.HTTPRequest.Context(), *params.Body.IP, *params.Body.Login, *params.Body.Password)
	if err != nil {
		return operations.NewAuthCheckInternalServerError().WithPayload(&operations.AuthCheckInternalServerErrorBody{
			Passed: false,
			Error:  err.Error(),
		})
	}
	fmt.Println(pass)
	return operations.NewAuthCheckOK().WithPayload(&operations.AuthCheckOKBody{Passed: pass})
}
