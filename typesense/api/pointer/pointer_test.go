package pointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointerValues(t *testing.T) {
	assert.NotNil(t, True())
	assert.Equal(t, true, *True())
	assert.NotNil(t, False())
	assert.Equal(t, false, *False())
	assert.NotNil(t, Int(10))
	assert.Equal(t, 10, *Int(10))
	assert.NotNil(t, String("abc"))
	assert.Equal(t, "abc", *String("abc"))

	var expectedFloat32 float32 = 9.5
	assert.NotNil(t, Float32(9.5))
	assert.Equal(t, expectedFloat32, *Float32(9.5))

	assert.NotNil(t, Float64(9.5))
	assert.Equal(t, 9.5, *Float64(9.5))

	v := struct{ field string }{field: "abc"}
	assert.NotNil(t, Interface(v))
	assert.Equal(t, v, *Interface(v))
}
