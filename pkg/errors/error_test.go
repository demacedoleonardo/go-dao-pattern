package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "client_side", E4xxCLIENTSIDE)
	assert.Equal(t, "unauthorized", E4xxUNAUTHORIZED)
	assert.Equal(t, "unprocessable_entity", E4xxUNPROCESSABLE)
	assert.Equal(t, "internal", E5xxINTERNAL)
}

func TestError_Error(t *testing.T) {
	t.Parallel()
	e := Error{
		Code:    E5xxINTERNAL,
		Message: "internal",
	}

	assert.Equal(t, "error: code=internal message=internal", e.Error())
}

func TestErrorCode(t *testing.T) {
	t.Parallel()
	assert.Equal(t, E5xxINTERNAL, ErrorCode(errors.New("default internal error")))
	assert.Equal(t, E4xxUNPROCESSABLE, ErrorCode(Errorf(E4xxUNPROCESSABLE, "error error error")))
	assert.Equal(t, "", ErrorCode(nil))
}

func TestIs_Error(t *testing.T) {
	t.Parallel()
	assert.True(t, Is(E4xxUNPROCESSABLE, Errorf(E4xxUNPROCESSABLE, "error error error")))
	assert.False(t, Is(E5xxINTERNAL, Errorf(E4xxUNPROCESSABLE, "error error error")))
	assert.False(t, Is(E5xxINTERNAL, nil))
	assert.False(t, Is(E5xxINTERNAL, errors.New("default false")))
}

func TestErrorMessage(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "internal error", ErrorMessage(errors.New("default internal error")))
	assert.Equal(t, "error error error", ErrorMessage(Errorf(E4xxUNPROCESSABLE, "error error error")))
	assert.Equal(t, "", ErrorMessage(nil))
}

func TestErrorStatus(t *testing.T) {
	t.Parallel()
	assert.Equal(t, http.StatusInternalServerError, ErrorStatus(Errorf(E5xxINTERNAL, "error error error")))
	assert.Equal(t, http.StatusUnprocessableEntity, ErrorStatus(Errorf(E4xxUNPROCESSABLE, "error error error")))
	assert.Equal(t, http.StatusNotFound, ErrorStatus(Errorf(E4xxNOTFOUND, "error error error")))
	assert.Equal(t, http.StatusUnauthorized, ErrorStatus(Errorf(E4xxUNAUTHORIZED, "error error error")))
	assert.Equal(t, http.StatusBadRequest, ErrorStatus(Errorf(E4xxCLIENTSIDE, "error error error")))
}
