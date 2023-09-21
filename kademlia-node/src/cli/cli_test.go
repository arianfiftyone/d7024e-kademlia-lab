package cli

import (
	"bytes"
	"testing"

	"github.com/arianfiftyone/src/kademlia"
	"github.com/stretchr/testify/assert"
)

func TestNewCli(t *testing.T) {
	// Create a test Kademlia instance.
	testKademlia := kademlia.NewKademlia("localhost", 9000, true, "10.0.0.1", 100000)

	// Create a new Cli instance.
	cli := NewCli(testKademlia)

	// Verify that the created Cli instance has the correct Kademlia instance.
	assert.Equal(t, testKademlia, cli.kademlia)
}

func TestLogMode(t *testing.T) {

	// Create a capturing buffer to capture the output.
	capturingBuffer := &bytes.Buffer{}
	output = capturingBuffer

	cli := &Cli{}

	// Verify that we are in LogMode
	mode := cli.LogMode(output)
	assert.Equal(t, LOG, mode)

	// Verify that the captured output contains the prompt message.
	outputString := capturingBuffer.String()
	assert.Contains(t, outputString, "You are in log mode, press `enter` to change to command mode")
}

func TestCommandMode(t *testing.T) {

	// Create a capturing buffer to capture the output.
	capturingBuffer := &bytes.Buffer{}
	output = capturingBuffer

	cli := &Cli{}

	// Verify that we are in CommandMode
	mode := cli.CommandMode(nil)
	assert.Equal(t, COMMAND, mode)

	// Verify that the captured output contains the expected prompt message.
	outputString := capturingBuffer.String()
	assert.Contains(t, outputString, "You are in command mode, enter `stop` to change to log mode")
}

func TestTrimWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  testing testy  ", "testing testy"}, // Leading and trailing spaces should be removed.
		{"   ", ""},                            // Only spaces should become an empty string.
		{"testy", "testy"},                     // No spaces, so the string should remain the same.
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := trimWhitespace(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestSplitInput(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"testing, testy ", []string{"testing,", "testy"}},                                  // Words separated by a comma and space.
		{"  The test is  being   tested", []string{"The", "test", "is", "being", "tested"}}, // Multiple spaces as separators.
		{"", []string{}}, // Empty input should result in an empty slice.
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := splitInput(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
