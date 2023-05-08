package common

import (
	"testing"

	"github.com/partyscript/bledger/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildtimeError(t *testing.T) {
	err := BLedgerBuildtimeError
	assert.Equal(t, err.Message, "unknown error occurred when starting the server")
	assert.Equal(t, err.Status, 503)
}

func TestBadRequestError(t *testing.T) {
	err := BLedgerBadRequestError
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.Message, "bad request")
	assert.Equal(t, err.Status, 400)
}

func TestBuildtimeWrappedError(t *testing.T) {
	err := WrapError("bad buildtime", BLedgerBuildtimeError, model.APIVersion("v1"))
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.APIVersion, model.APIVersion("v1"))
	assert.Equal(t, err.InternalErrMsg, "bad buildtime")
	assert.Equal(t, err.Error.(StandardSentinelError).Message, "unknown error occurred when starting the server")
	assert.Equal(t, err.Error.(StandardSentinelError).Status, 503)
}

func TestBadRequestWrappedError(t *testing.T) {
	err := WrapError("bad req", BLedgerBadRequestError, model.APIVersion("v1"))
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.APIVersion, model.APIVersion("v1"))
	assert.Equal(t, err.InternalErrMsg, "bad req")
	assert.Equal(t, err.Error.(StandardSentinelError).Message, "bad request")
	assert.Equal(t, err.Error.(StandardSentinelError).Status, 400)

}

func TestBuildtimeWrappedAPIError(t *testing.T) {
	status, err := WrapAPIError("bad buildtime", BLedgerBuildtimeError, model.APIVersion("v1"))
	assert.NotEqual(t, err, nil)
	assert.NotEqual(t, status, 0)
	assert.Equal(t, status, 503)
	assert.Equal(t, err.APIVersion, model.APIVersion("v1"))
	assert.Equal(t, err.InternalErrMsg, "bad buildtime")
	assert.Equal(t, err.Error.(StandardSentinelError).Message, "unknown error occurred when starting the server")
	assert.Equal(t, err.Error.(StandardSentinelError).Status, 503)
}

func TestBadRequestWrappedAPIError(t *testing.T) {
	status, err := WrapAPIError("bad req", BLedgerBadRequestError, model.APIVersion("v1"))
	assert.NotEqual(t, err, nil)
	assert.NotEqual(t, status, 0)
	assert.Equal(t, status, 400)
	assert.Equal(t, err.APIVersion, model.APIVersion("v1"))
	assert.Equal(t, err.InternalErrMsg, "bad req")
	assert.Equal(t, err.Error.(StandardSentinelError).Message, "bad request")
	assert.Equal(t, err.Error.(StandardSentinelError).Status, 400)

}
