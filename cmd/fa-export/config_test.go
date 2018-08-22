package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ()

func TestLoadConfigNone(t *testing.T) {
	baseDir = "test"
	configName = "/nonexistant"

	c, err := LoadConfig(configName)
	assert.Nil(t, c, "nil")
	assert.Error(t, err, "should be in error")
}

func TestLoadConfigBad(t *testing.T) {
	baseDir = "test"
	configName = "bad.toml"

	c, err := LoadConfig(configName)
	assert.Nil(t, c, "nil value")
	assert.Error(t, err, "should be in error")
}

func TestLoadConfigPerms(t *testing.T) {
	baseDir = "test"
	configName = "config.toml"

	file := filepath.Join(baseDir, configName)
	err := os.Chmod(file, 0000)
	assert.NoError(t, err, "should be fine")

	c, err := LoadConfig(configName)
	assert.Nil(t, c, "nil value")
	assert.Error(t, err, "should be in error")

	err = os.Chmod(file, 0644)
	assert.NoError(t, err, "should be fine")
}

func TestLoadConfigGood(t *testing.T) {
	baseDir = "test"
	configName = "config.toml"

	c, err := LoadConfig(configName)
	assert.NotNil(t, c, "not nil")
	assert.NoError(t, err, "should be fine")

	// Check values
	assert.EqualValues(t, testCnf, c, "should be equal")
}

func TestLoadConfigGoodVerbose(t *testing.T) {
	baseDir = "test"
	configName = "config.toml"
	fVerbose = true

	c, err := LoadConfig(configName)
	assert.NotNil(t, c, "not nil")
	assert.NoError(t, err, "should be fine")

	// Check values
	assert.EqualValues(t, testCnf, c, "should be equal")
	fVerbose = false
}
