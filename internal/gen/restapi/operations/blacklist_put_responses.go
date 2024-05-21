// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"OTUS_hws/Anti-BruteForce/internal/gen/models"
)

// BlacklistPutOKCode is the HTTP code returned for type BlacklistPutOK
const BlacklistPutOKCode int = 200

/*BlacklistPutOK The IP has been successfully added to blacklist

swagger:response blacklistPutOK
*/
type BlacklistPutOK struct {

	/*
	  In: Body
	*/
	Payload *models.Status200 `json:"body,omitempty"`
}

// NewBlacklistPutOK creates BlacklistPutOK with default headers values
func NewBlacklistPutOK() *BlacklistPutOK {

	return &BlacklistPutOK{}
}

// WithPayload adds the payload to the blacklist put o k response
func (o *BlacklistPutOK) WithPayload(payload *models.Status200) *BlacklistPutOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the blacklist put o k response
func (o *BlacklistPutOK) SetPayload(payload *models.Status200) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BlacklistPutOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// BlacklistPutInternalServerErrorCode is the HTTP code returned for type BlacklistPutInternalServerError
const BlacklistPutInternalServerErrorCode int = 500

/*BlacklistPutInternalServerError Internal Server Error

swagger:response blacklistPutInternalServerError
*/
type BlacklistPutInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error500 `json:"body,omitempty"`
}

// NewBlacklistPutInternalServerError creates BlacklistPutInternalServerError with default headers values
func NewBlacklistPutInternalServerError() *BlacklistPutInternalServerError {

	return &BlacklistPutInternalServerError{}
}

// WithPayload adds the payload to the blacklist put internal server error response
func (o *BlacklistPutInternalServerError) WithPayload(payload *models.Error500) *BlacklistPutInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the blacklist put internal server error response
func (o *BlacklistPutInternalServerError) SetPayload(payload *models.Error500) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BlacklistPutInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
