package flightaware

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFAClient_Verbose0(t *testing.T) {
	var buf bytes.Buffer
	var deflog = log.New(&buf, "test: ", log.Lshortfile)

	c := NewClient(rc)
	require.NotNil(t, c)

	c.SetLog(deflog)
	c.verbose("foo")
	assert.Empty(t, buf.String())
}

func TestFAClient_Verbose1(t *testing.T) {
	var buf bytes.Buffer
	var deflog = log.New(&buf, "test: ", log.Lshortfile)

	c := NewClient(rc)
	require.NotNil(t, c)

	c.SetLog(deflog)
	c.SetLevel(1)
	c.verbose("foo")
	assert.NotEmpty(t, buf.String())
	assert.True(t, strings.HasPrefix(buf.String(), "test:"))
	assert.True(t, strings.HasSuffix(buf.String(), "foo\n"))
}

func TestFAClient_Verbose2(t *testing.T) {
	var buf bytes.Buffer
	var deflog = log.New(&buf, "test: ", log.Lshortfile)

	c := NewClient(rc)
	require.NotNil(t, c)

	c.SetLog(deflog)
	c.SetLevel(2)
	c.verbose("foo")
	assert.NotEmpty(t, buf.String())
	assert.True(t, strings.HasPrefix(buf.String(), "test:"))
	assert.True(t, strings.HasSuffix(buf.String(), "foo\n"))
}

func TestFAClient_Debug0(t *testing.T) {
	var buf bytes.Buffer
	var deflog = log.New(&buf, "test: ", log.Lshortfile)

	c := NewClient(rc)
	require.NotNil(t, c)

	c.SetLog(deflog)
	c.SetLevel(0)
	c.debug("foo")
	assert.Empty(t, buf.String())
}

func TestFAClient_Debug1(t *testing.T) {
	var buf bytes.Buffer
	var deflog = log.New(&buf, "test: ", log.Lshortfile)

	c := NewClient(rc)
	require.NotNil(t, c)

	c.SetLog(deflog)
	c.SetLevel(1)
	c.debug("foo")
	assert.Empty(t, buf.String())
}

func TestFAClient_Debug2(t *testing.T) {
	var buf bytes.Buffer
	var deflog = log.New(&buf, "test: ", log.Lshortfile)

	c := NewClient(rc)
	require.NotNil(t, c)

	c.SetLog(deflog)
	c.SetLevel(2)
	c.debug("foo")
	assert.NotEmpty(t, buf.String())
	assert.True(t, strings.HasPrefix(buf.String(), "test:"))
	assert.True(t, strings.HasSuffix(buf.String(), "foo\n"))
}
