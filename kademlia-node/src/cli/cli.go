package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/arianfiftyone/src/kademlia"

	"github.com/arianfiftyone/src/logger"
)

// `input` is a variable that represents standard input (os.Stdin)
var input *os.File = os.Stdin

// `output` is a variable that represents standard output (os.Stdout)
var output io.Writer = os.Stdout

type Mode string

const (
	COMMAND Mode = "COMMAND"
	LOG     Mode = "LOG"
)

type Cli struct {
	kademlia kademlia.Kademlia
}

func NewCli(kademlia kademlia.Kademlia) *Cli {
	return &Cli{
		kademlia: kademlia,
	}
}

func (cli *Cli) Clear() {
	cmd := exec.Command("clear") // Linux tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (cli *Cli) showOldLogs() {
	cli.Clear()
	oldLogs := logger.GetOldLogs()
	for _, log := range oldLogs {
		fmt.Println(log)
	}
}

func (cli *Cli) StartCli(out io.Writer) {
	mode := LOG
	for {
		if mode == LOG {
			mode = cli.LogMode(out)
		} else {
			mode = cli.CommandMode(out)
		}
	}
}

func (cli *Cli) LogMode(out io.Writer) Mode {

	// Disable echoing to the terminal
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	fmt.Fprintln(output, "You are in log mode, press `enter` to change to command mode")

	// Start a goroutine for continuous log display
	stopLogChannel := make(chan bool)
	go func() {
	LogLoop:
		for {
			select {
			case <-stopLogChannel:
				break LogLoop

			// Read new log entries
			case <-time.After(time.Millisecond * 100):
				logEntry, err := logger.ReadNewLog()
				if err == nil {
					fmt.Println(logEntry)
				}
			}
		}
	}()

	inputBuffer := make([]byte, 1)
	n, _ := os.Stdin.Read(inputBuffer)

	stopLogChannel <- true

	// Toggle between `COMMAND` and `LOG` based on ENTER keypress
	if string(inputBuffer[:n]) == "\n" {
		cli.Clear()
		return COMMAND
	} else {
		return LOG
	}
}

func (cli *Cli) CommandMode(out io.Writer) Mode {

	// Configure the terminal for echo mode
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	fmt.Fprintln(output, "You are in command mode, enter `stop` to change to log mode")

	fmt.Print("$ ")

	// Read user input from stdin
	reader := bufio.NewReader(input)
	input, _ := reader.ReadString('\n')

	trimmedInput := trimWhitespace(input)
	if trimmedInput == "stop" {
		cli.showOldLogs()
		return LOG
	}
	if trimmedInput == "" {
		return COMMAND
	} else {
		// Split input into individual commands
		commands := splitInput(trimmedInput)

		cli.HandleCommands(out, cli.kademlia, commands)
		return COMMAND
	}
}

// `trimWhitespace` removes leading and trailing whitespace from a string.
func trimWhitespace(input string) string {
	return strings.TrimSpace(input)
}

// `splitInput` splits a string into individual words separated with a whitespace.
func splitInput(input string) []string {
	return strings.Fields(input)
}
