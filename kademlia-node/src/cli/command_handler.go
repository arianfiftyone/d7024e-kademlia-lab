package cli

import (
	"fmt"
	"io"
	"log"
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
func (cli *Cli) HandleCommands(output io.Writer, kademliaInstance *kademlia.Kademlia, commands []string) {

	numArgs := len(commands)
	command := strings.ToLower(commands[0])

	switch command {

	case "put", "p":
		if numArgs == 2 {
			Put(*kademliaInstance, commands[1])
		} else {
			fmt.Fprintln(output, noArgsError)
		}

	case "get", "g":
		if numArgs == 2 {
			Get(*kademliaInstance, kademlia.GetKeyRepresentationOfKademliaId(kademlia.NewKademliaID(commands[1])))
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
			hexRepresentation := *kademliaInstance.KademliaNode.GetRoutingTable().Me.ID
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

func Put(kademlia kademlia.Kademlia, content string) {
	key, err := kademlia.Store(content)

	if err != nil {
		log.Printf("Error when storing content: %v\n", err)
	} else {
		fmt.Printf("Got hash: %s\n", key.GetHashString())
	}

}

func Get(kademlia kademlia.Kademlia, key *kademlia.Key) {
	_, value, err := kademlia.LookupData(key)

	if err != nil {
		log.Printf("Error when looking up data %v\n", err)
		return
	}

	if value != "" {
		fmt.Printf("Got content: %s\n", value)
		return
	}
}

func Kill() {
	fmt.Println("Terminating node and exiting...")
	exit(1)
}

func Help(output io.Writer) {
	text := HelpPrompt()
	fmt.Fprintln(output, text)
}
