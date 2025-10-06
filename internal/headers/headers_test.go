package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
func TestHeaders_MultipleValues(t *testing.T) {
	headers := NewHeaders()

	data := []byte(
		"Host: localhost:42069\r\nHost: localhost:42069\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069,localhost:42069", headers["host"])
	assert.False(t, done)

}
func TestHeadersMixedCaseKey(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-CuStOm-KeY: SomeValue\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.True(t, done)
	// La clé doit être stockée en minuscules
	assert.Equal(t, "SomeValue", headers["x-custom-key"])
}

func TestInvalidCharacterInHeaderKey(t *testing.T) {
	headers := NewHeaders()
	data := []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
