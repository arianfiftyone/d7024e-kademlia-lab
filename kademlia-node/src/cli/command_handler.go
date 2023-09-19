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
func HandleCommands(output io.Writer, kademlia *kademlia.Kademlia, commands []string) {

	numArgs := len(commands)
	command := strings.ToLower(commands[0])

	switch command {

	case "put", "p":
		if numArgs == 2 {
			// Put(*kademlia, commands[1])
			fmt.Printf("PUT is not implemented yet: %s\n", commands[1])
		} else {
			fmt.Fprintln(output, noArgsError)
		}

	case "get", "g":
		if numArgs == 2 {
			// Get(*kademlia, commands[1])
			fmt.Printf("GET is not implemented yet: %s\n", commands[1])
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
			hexRepresentation := *kademlia.KademliaNode.RoutingTable.Me.ID
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

	default:
		fmt.Fprintln(output, commandError)
	}

}

func Put(kademlia kademlia.Kademlia, content string) {
	hash, err := kademlia.Store(content) // nothing to return for now

	if err != nil {
		log.Printf("Error when storing content: %v\n", err)
	} else {
		fmt.Printf("Got hash: %s\n", hash)
	}

}

func Get(kademlia kademlia.Kademlia, key *kademlia.Key) {
	content, err := kademlia.LookupData(key) // nothing to return for now

	if err != nil {
		log.Printf("Error when looking up data %v\n", err)
	} else {
		fmt.Printf("Got content: %s\n", content)
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
