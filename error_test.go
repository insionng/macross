package macross

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpError(t *testing.T) {
	e := NewHTTPError(StatusNotFound)
	assert.Equal(t, StatusNotFound, e.StatusCode())
	assert.Equal(t, StatusText(StatusNotFound), e.Error())

	e = NewHTTPError(StatusNotFound, "abc")
	assert.Equal(t, StatusNotFound, e.StatusCode())
	assert.Equal(t, "abc", e.Error())

	s, _ := json.Marshal(e)
	assert.Equal(t, `{"Status":404,"Message":"abc"}`, string(s))
}
