package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/arianfiftyone/src/kademlia"
	"github.com/stretchr/testify/assert"
)

func testCommand(command string) string {
	output = bytes.NewBuffer(nil)

	HandleCommands(output, kademlia.NewKademlia("localhost", 9000, true, "10.0.0.1", 100000), []string{command})
	return trimNewlineFromWriterOutput(output)
}

func trimNewlineFromWriterOutput(output io.Writer) string {
	s := output.(*bytes.Buffer).String()
	return strings.TrimSuffix(s, "\n")
}

func TestGetCommand(t *testing.T) {
	assert.Equal(t, noArgsError, testCommand("get"))
}

func TestAbbreviatedGetCommand(t *testing.T) {
	assert.Equal(t, noArgsError, testCommand("g"))
}

func TestPutCommand(t *testing.T) {
	assert.Equal(t, noArgsError, testCommand("put"))
}

func TestAbbreviatedPutCommand(t *testing.T) {
	assert.Equal(t, noArgsError, testCommand("p"))
}

func TestHelpCommand(t *testing.T) {
	content := HelpPrompt()
	assert.Equal(t, content, testCommand("help"))
}

func TestAbbreviatedHelpCommand(t *testing.T) {
	content := HelpPrompt()
	assert.Equal(t, content, testCommand("h"))
}

func TestKademliaIDCommand(t *testing.T) {
	assert.Equal(t, "ffffffff00000000000000000000000000000000", testCommand("kademliaid"))
}

func TestAbbreviatedKademliaIDCommand(t *testing.T) {
	assert.Equal(t, "ffffffff00000000000000000000000000000000", testCommand("kid"))
}

func TestKillCommand(t *testing.T) {
	assert.Equal(t, 1, exitCli("kill"))
}

func TestAbbreviatedKillCommand(t *testing.T) {
	assert.Equal(t, 1, exitCli("k"))
}

// `exitCLI` purpose is for testing cases(`kill`) where the function exits with `os.Exit()`
// and is slightly modified from the source:
// https://stackoverflow.com/questions/40615641/testing-os-exit-scenarios-in-go-with-coverage-information-c overalls-io-goverall/40801733#40801733
func exitCli(e string) int {
	var got int
	// Save current function and restore at the end:
	oldExit := exit

	// Replace exit with a custom function that captures the exit code.
	exit = func(code int) {
		got = code
	}

	defer func() {
		// Restore the original exit function when exiting the function.
		exit = oldExit
	}()

	HandleCommands(output, nil, []string{e})
	return got
}
