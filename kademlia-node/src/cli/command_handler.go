package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arianfiftyone/src/kademlia"
)

const (
	noArgsError       = "Please provide a correct ARGUMENT!"
	commandError      = "Please provide a correct COMMAND!"
	fileNotFoundError = "Could find not and open the file: "
)

var (
	exit = os.Exit
)

// HandleCommand handles CLI commands.
// `output` is the io.Writer used for data output.
// `kademlia` represents the Kademlia node associated with this CLI
// `commands` is a slice of program commands.
func (cli *Cli) HandleCommands(output io.Writer, kademliaInstance kademlia.Kademlia, commands []string) {

	numArgs := len(commands)
	command := strings.ToLower(commands[0])

	switch command {

	case "put", "p":
		if numArgs == 2 {
			key, err := Put(kademliaInstance, commands[1])
			if err != nil {
				customErr := fmt.Errorf("error when storing content: %s", err.Error())
				fmt.Fprintln(output, customErr)
			} else {
				fmt.Fprintln(output, "Got hash: "+key)

			}

		} else {
			fmt.Fprintln(output, noArgsError)
		}

	case "get", "g":
		if numArgs == 2 {
			content, err := Get(kademliaInstance, kademlia.GetKeyRepresentationOfKademliaId(kademlia.NewKademliaID(commands[1])))
			if err != nil {
				customErr := fmt.Errorf("error when looking up data %s", err.Error())
				fmt.Fprintln(output, customErr)
			} else {
				fmt.Fprintln(output, "Got content: "+content)

			}
		} else {
			fmt.Fprintln(output, noArgsError)
		}

	case "kill", "k":
		if numArgs == 1 {
			Kill()
		} else {
			fmt.Fprintln(output, noArgsError)
		}

	case "kademliaid", "kid":
		if numArgs == 1 {

			kademliaNode := *kademliaInstance.GetKademliaNode()
			hexRepresentation := kademliaNode.GetRoutingTable().Me.ID
			hexStringRepresentation := hexRepresentation.String()

			fmt.Fprintln(output, hexStringRepresentation)
		} else {
			fmt.Fprintln(output, noArgsError)
		}

	case "help", "h":
		if numArgs == 1 {
			Help(output)
		} else {
			fmt.Fprintln(output, noArgsError)
		}

	case "clear", "c":
		if numArgs == 1 {
			cli.Clear()
		} else {
			fmt.Fprintln(output, noArgsError)
		}

	default:
		fmt.Fprintln(output, commandError)
	}

}

func Put(kademlia kademlia.Kademlia, content string) (string, error) {
	key, err := kademlia.Store(content)

	if err != nil {
		return "", err
	} else {
		return key.GetHashString(), nil
	}

}

func Get(kademlia kademlia.Kademlia, key *kademlia.Key) (string, error) {
	_, value, err := kademlia.LookupData(key)

	if err != nil {
		return "", err
	}

	if value != "" {
		return value, nil
	}

	return "", nil
}

func Kill() {
	fmt.Println("Terminating node and exiting...")
	exit(1)
}

func Help(output io.Writer) {
	text := HelpPrompt()
	fmt.Fprintln(output, text)
}
