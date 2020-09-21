package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCmdInput_shouldReturnCommand(t *testing.T) {
	testParseCmdInput := func(inputString string, expected Command) {
		res, err := ParseCmdInput(inputString)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	}

	testParseCmdInput("command", Command{Command: "command", Args: []string{}})
	testParseCmdInput("command args1 args2", Command{Command: "command", Args: []string{"args1", "args2"}})
	testParseCmdInput("\"spaced command\" \"spaced args1\" args2", Command{Command: "spaced command", Args: []string{"spaced args1", "args2"}})
	testParseCmdInput("\"spaced command\" \"spaced missing close quote args1 args2", Command{Command: "spaced command", Args: []string{"spaced missing close quote args1 args2"}})
	testParseCmdInput("\"spaced command\" spaced\" joined string \"args2", Command{Command: "spaced command", Args: []string{"spaced joined string args2"}})
}

func TestParseCmdInput_shouldReturnError_givenNoCommand(t *testing.T) {
	testParseCmdInput := func(inputString string) {
		res, err := ParseCmdInput(inputString)
		assert.Error(t, err)
		assert.Equal(t, "no command given", err.Error())
		assert.Equal(t, Command{}, res)
	}

	testParseCmdInput("")
}
