package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpPrompt(t *testing.T) {
	expectedHelpText := `
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
								[ KADEMLIA COMMAND LINE INTERFACE ]
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

COMMANDS:
	command [command options] [args...]

VERSION:
	v1.0
	
COMMANDS:
	get, g <hash>      		Takes the hash and outputs the contents of the object and the node it was retrieved from, if it could be downloaded
	put, p <content>      		Takes the content of the file you are uploading and outputs the hash of the object, if content could be uploaded
	kill, k      			Kills the node
	kademliaid, kid 		Get id associated with the node	 
	help, h      			Output this help prompt
	clear, c			Clear the terminal

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
							[ MADE BY: arianfiftyone, MrDweller & asta987 ]
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
`

	result := HelpPrompt()
	assert.Equal(t, expectedHelpText, result)
}
