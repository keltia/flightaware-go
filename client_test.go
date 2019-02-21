package flightaware

import (
	"bytes"
	"log"
	"testing"
	"time"

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

func TestFAClient_SetTimer(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)

	c1 := c.SetTimer(int64(3600))
	require.Equal(t, c, c1)
}

func TestFAClient_SetFeed(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)

	err := c.SetFeed("", []time.Time{})
	require.NoError(t, err)
}

func TestFAClient_SetFeed2(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)

	futur := time.Now().Add(10 *time.Minute)
	err := c.SetFeed("pitr", []time.Time{futur})
	require.Error(t, err)
}

func TestFAClient_SetFeed3(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)

	past := time.Now().Add(-10 * time.Minute)
	err := c.SetFeed("pitr", []time.Time{past})
	require.NoError(t, err)
	assert.Equal(t, past.Unix(), c.RangeT[0])
}

func TestFAClient_SetFeed4(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)

	beg := time.Now().Add(-10 * time.Minute)
	end := time.Now().Add(-5 * time.Minute)

	err := c.SetFeed("range", []time.Time{beg, end})
	require.NoError(t, err)
	assert.Equal(t, beg.Unix(), c.RangeT[0])
	assert.Equal(t, end.Unix(), c.RangeT[1])

}

func TestFAClient_Version(t *testing.T) {
	c := NewClient(rc)
	require.NotNil(t, c)
	require.Equal(t, FAVersion, c.Version())
}
