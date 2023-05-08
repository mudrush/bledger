package common

import (
	"net/http"

	"github.com/partyscript/bledger/internal/model"
)

var (
	// BLedgerBuildtimeError is an error used when there are buildtime issues
	BLedgerBuildtimeError = StandardSentinelError{
		Status:  http.StatusServiceUnavailable,
		Message: "unknown error occurred when starting the server",
	}

	// BLedgerBadRequestError is an error used to show a request was invalid
	BLedgerBadRequestError = StandardSentinelError{
		Status:  http.StatusBadRequest,
		Message: "bad request",
	}

	// BLedgerIdempotencyError is an error used to show a request was invalid
	BLedgerIdempotencyError = StandardSentinelError{
		Status:  http.StatusForbidden,
		Message: "idempotency key is invalid",
	}

	// BLedgerNotFoundError is an error used to show a request was invalid
	BLedgerNotFoundError = StandardSentinelError{
		Status:  http.StatusNotFound,
		Message: "not found",
	}

	// BLedgerInternalError is an error used to show a request was invalid
	BLedgerInternalError = StandardSentinelError{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}
)

// APIError is an interface for accessing and implementing a client error
type APIError interface {
	APIError() (int, string)
}

// StandardSentinelError is a struct for returning embedded runtime errors
type StandardSentinelError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// WrapError is used to create a public facing wrapped errors that
// includes the `sentinelError` which is an internal or service application errors
func WrapError(errMsg string, sentinel StandardSentinelError, av model.APIVersion) model.StandardErrorResponse {
	return model.StandardErrorResponse{InternalErrMsg: errMsg, Error: sentinel, APIVersion: av}
}

// WrapAPIError is used to create a public facing wrapped errors that
// includes the `sentinelError` which is an internal or service application errors
func WrapAPIError(errMsg string, sentinel StandardSentinelError, av model.APIVersion) (int, model.StandardErrorResponse) {
	return sentinel.Status, model.StandardErrorResponse{InternalErrMsg: errMsg, Error: sentinel, APIVersion: av}
}
