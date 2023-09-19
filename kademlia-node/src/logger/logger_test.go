package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	// Reset logger
	logger = nil

	Log("Test Log 1")
	assert.Equal(t, []string{"Test Log 1"}, logger.newLogs)

	Log("Test Log 2")
	assert.Equal(t, []string{"Test Log 1", "Test Log 2"}, logger.newLogs)
}

func TestReadNewLog(t *testing.T) {
	// Reset logger
	logger = nil

	_, err := ReadNewLog()
	assert.Error(t, err, "No new logs available")

	Log("Test Log 1")
	log, err := ReadNewLog()
	assert.NoError(t, err)
	assert.Equal(t, "Test Log 1", log)

	Log("Test Log 2")
	log, err = ReadNewLog()
	assert.NoError(t, err)
	assert.Equal(t, "Test Log 2", log)

	_, err = ReadNewLog()
	assert.Error(t, err, "No new logs available")
}

func TestGetOldLogs(t *testing.T) {
	// Reset logger
	logger = nil

	oldLogs := GetOldLogs()
	assert.Empty(t, oldLogs) // `oldLogs` should be empty until `ReadNewLog()` is called.

	Log("Test Log 1")
	_, err := ReadNewLog()
	assert.NoError(t, err)

	oldLogs = GetOldLogs()
	assert.Equal(t, []string{"Test Log 1"}, oldLogs)
}
