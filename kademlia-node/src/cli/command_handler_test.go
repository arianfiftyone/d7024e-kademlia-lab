package cli

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/arianfiftyone/src/kademlia"
	"github.com/stretchr/testify/assert"
)

// Define a function to create a test Kademlia instance
func createTestKademlia() *kademlia.Kademlia {
	kademlia := kademlia.NewKademlia("localhost", 9000, true, "10.0.0.1", 100000)
	return kademlia
}

func (cli *Cli) testCommand(command string) string {
	output = bytes.NewBuffer(nil)

	cli.HandleCommands(output, createTestKademlia(), []string{command})
	return trimNewlineFromWriterOutput(output)
}

func trimNewlineFromWriterOutput(output io.Writer) string {
	s := output.(*bytes.Buffer).String()
	return strings.TrimSuffix(s, "\n")
}

func TestGetCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, noArgsError, cli.testCommand("get"))
}

func TestAbbreviatedGetCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, noArgsError, cli.testCommand("g"))
}

func TestPutCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, noArgsError, cli.testCommand("put"))
}

func TestAbbreviatedPutCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, noArgsError, cli.testCommand("p"))
}

func TestHelpCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	content := HelpPrompt()
	assert.Equal(t, content, cli.testCommand("help"))
}

func TestAbbreviatedHelpCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	content := HelpPrompt()
	assert.Equal(t, content, cli.testCommand("h"))
}

func TestKademliaIDCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, "ffffffff00000000000000000000000000000000", cli.testCommand("kademliaid"))
}

func TestAbbreviatedKademliaIDCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, "ffffffff00000000000000000000000000000000", cli.testCommand("kid"))
}

func TestKillCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, 1, cli.exitCli("kill"))
}

func TestAbbreviatedKillCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, 1, cli.exitCli("k"))
}

func TestClearCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, "", cli.testCommand("clear"))
}

func TestAbbreviatedClearCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)
	assert.Equal(t, "", cli.testCommand("c"))
}

// `exitCLI` purpose is for testing cases(`kill`) where the function exits with `os.Exit()`
// and is slightly modified from the source:
// https://stackoverflow.com/questions/40615641/testing-os-exit-scenarios-in-go-with-coverage-information-c overalls-io-goverall/40801733#40801733
func (cli *Cli) exitCli(e string) int {
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

	cli.HandleCommands(output, nil, []string{e})
	return got
}

// Not sure if tested correctly
func TestClear(t *testing.T) {
	// Backup the original os.Stdout and restore it at the end of the test.
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	// Create a capturing buffer to capture the output.
	capturingBuffer := &bytes.Buffer{}
	log.SetOutput(capturingBuffer)

	cli := &Cli{}

	cli.Clear()

	// Define the expected clear command sequence for your system.
	expectedClear := ""

	// Verify that the captured output contains the expected clear command sequence.
	output := capturingBuffer.String()
	assert.Equal(t, output, expectedClear)
}
