package cli

import "errors"

// Command is the struct used to define the command and arguments
type Command struct {
	Command string
	Args    []string
}

// ParseCmdInput breaks down inputs to command and parameters
func ParseCmdInput(input string) (Command, error) {
	var quotation bool
	var cur string
	parts := []string{}

	for _, char := range input {
		switch char {
		case '"':
			quotation = !quotation
		case ' ':
			if !quotation && cur != "" {
				parts = append(parts, cur)
				cur = ""
			} else {
				cur += string(char)
			}
		default:
			cur += string(char)
		}
	}
	if cur != "" {
		parts = append(parts, cur)
	}

	if len(parts) == 0 {
		return Command{}, errors.New("no command given")
	}

	return Command{Command: parts[0], Args: parts[1:]}, nil
}
