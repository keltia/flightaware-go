package flightaware

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var rc = Config{}

func TestNewClient(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)
	require.IsType(t, (*FAClient)(nil), c)
	assert.Equal(t, rc, c.Host)
}

func TestFAClient_AddHandler(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)
	require.IsType(t, (*FAClient)(nil), c)

	var fn = func(f []byte) {}

	c.AddHandler(fn)
	require.NotNil(t, c.FeedOne)
}

func TestFAClient_SetLog(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)
	require.IsType(t, (*FAClient)(nil), c)

	ol := c.Log

	var buf bytes.Buffer

	tl := log.New(&buf, "test", log.LstdFlags)
	c.SetLog(tl)
	require.NotEqual(t, ol, c.Log)
}

func TestFAClient_SetLevel(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)

	c.SetLevel(0)
	require.Equal(t, 0, c.level)
}

func TestFAClient_SetLevel1(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)

	c.SetLevel(1)
	require.Equal(t, 1, c.level)
}

func TestFAClient_SetLevel2(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)

	c.SetLevel(2)
	require.Equal(t, 2, c.level)
}

func TestFAClient_Close(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)

	require.Error(t, c.Close())
}

func TestFAClient_Version(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)
	require.Equal(t, FAVersion, c.Version())
}
