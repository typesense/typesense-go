package typesense

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpError(t *testing.T) {
	err := &httpError{status: 200, body: []byte("error message body")}
	assert.Equal(t, "status: 200 response: error message body", err.Error())
}
