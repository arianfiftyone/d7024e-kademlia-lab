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

type KademliaMock struct {
	DataStore *kademlia.DataStore
}

func (KademliaMock *KademliaMock) Start() {}
func (KademliaMock *KademliaMock) Join()  {}
func (KademliaMock *KademliaMock) Store(content string) (*kademlia.Key, error) {
	return kademlia.HashToKey(content), nil
}

func (KademliaMock *KademliaMock) GetKademliaNode() *kademlia.KademliaNode {
	return nil
}

func (KademliaMock *KademliaMock) FirstSetContainsAllContactsOfSecondSet(first []kademlia.Contact, second []kademlia.Contact) bool {
	return false
}

func (KademliaMock *KademliaMock) LookupContact(targetId *kademlia.KademliaID) ([]kademlia.Contact, error) {
	return nil, nil
}

func (KademliaMock *KademliaMock) LookupData(key *kademlia.Key) ([]kademlia.Contact, string, error) {
	content, err := KademliaMock.DataStore.Get(key)
	return nil, content, err
}

// Creates a test Kademlia instance
func createTestKademlia() *kademlia.KademliaImplementation {
	kademlia := kademlia.NewKademlia("localhost", 9000, true, "10.0.0.1", 100000)
	return kademlia
}

func (cli *Cli) testCommand(commands []string) string {
	output = bytes.NewBuffer(nil)

	cli.HandleCommands(output, cli.kademlia, commands)
	return trimNewlineFromWriterOutput(output)
}

func trimNewlineFromWriterOutput(output io.Writer) string {
	s := output.(*bytes.Buffer).String()
	return strings.TrimSuffix(s, "\n")
}

func TestGet(t *testing.T) {
	content := "kademlia"
	key := kademlia.HashToKey(content)

	dataStore := kademlia.NewDataStore()
	dataStore.Insert(key, content)

	cli := NewCli(&KademliaMock{
		DataStore: &dataStore,
	})

	command := []string{
		"get",
		key.GetHashString(),
	}

	output := cli.testCommand(command)
	assert.Equal(t, "Got content: "+content, output)
}

func TestGetCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)

	command := []string{
		"get",
	}

	output := cli.testCommand(command)
	assert.Equal(t, noArgsError, output)
}

func TestAbbreviatedGetCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()
	cli := NewCli(kademliaInstance)

	command := []string{
		"g",
	}

	output := cli.testCommand(command)
	assert.Equal(t, noArgsError, output)
}

func TestPut(t *testing.T) {
	value := "kademlia"
	key := kademlia.HashToKey(value)

	cli := NewCli(&KademliaMock{})
	command := []string{
		"put",
		value,
	}

	output := cli.testCommand(command)

	assert.Equal(t, "Got hash: "+key.GetHashString(), output)
}

func TestPutCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"put",
	}

	output := cli.testCommand(command)
	assert.Equal(t, noArgsError, output)
}

func TestAbbreviatedPutCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"p",
	}

	output := cli.testCommand(command)
	assert.Equal(t, noArgsError, output)
}

func TestKillCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"kill",
	}

	output := cli.exitCli(command)
	assert.Equal(t, 1, output)
}

func TestAbbreviatedKillCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"k",
	}

	output := cli.exitCli(command)
	assert.Equal(t, 1, output)
}

func TestKillCommandError(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"kill",
		"me",
	}

	output := cli.testCommand(command)
	assert.Equal(t, noArgsError, output)

}

func TestKademliaIDCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"kademliaid",
	}

	output := cli.testCommand(command)
	assert.Equal(t, "ffffffff00000000000000000000000000000000", output)
}

func TestAbbreviatedKademliaIDCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"kid",
	}

	output := cli.testCommand(command)
	assert.Equal(t, "ffffffff00000000000000000000000000000000", output)
}

func TestKademliaIDError(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"kademliaid",
		"me",
	}

	output := cli.testCommand(command)
	assert.Equal(t, noArgsError, output)
}

func TestHelpCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"help",
	}

	content := HelpPrompt()

	output := cli.testCommand(command)
	assert.Equal(t, content, output)
}

func TestAbbreviatedHelpCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"h",
	}

	content := HelpPrompt()

	output := cli.testCommand(command)
	assert.Equal(t, content, output)
}

func TestHelpCommandError(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"help",
		"me",
	}

	output := cli.testCommand(command)
	assert.Equal(t, noArgsError, output)
}

func TestClearCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"clear",
	}

	output := cli.testCommand(command)
	assert.Equal(t, "", output)
}

func TestAbbreviatedClearCommand(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"c",
	}

	output := cli.testCommand(command)
	assert.Equal(t, "", output)
}

func TestHelpClearError(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"clear",
		"me",
	}

	output := cli.testCommand(command)
	assert.Equal(t, noArgsError, output)
}

// `exitCLI` purpose is for testing cases(`kill`) where the function exits with `os.Exit()`
// and is slightly modified from the source:
// https://stackoverflow.com/questions/40615641/testing-os-exit-scenarios-in-go-with-coverage-information-c overalls-io-goverall/40801733#40801733
func (cli *Cli) exitCli(command []string) int {
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

	cli.HandleCommands(output, nil, command)
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

func TestCommandError(t *testing.T) {
	kademliaInstance := createTestKademlia()

	cli := NewCli(kademliaInstance)
	command := []string{
		"not correct input",
	}

	output := cli.testCommand(command)
	assert.Equal(t, commandError, output)
}
