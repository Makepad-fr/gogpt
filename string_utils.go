package gogpt

import "strings"

// isEmpty checks if the given string is empty or not by trimming it using strings.TrimSpace
func isEmpty(input string) bool {
	return len(strings.TrimSpace(input)) == 0
}

// isEndOfEventStream checks if the given input is the end of the event stream by the conversation API
func isEndOfEventStream(input string) bool {
	return strings.HasPrefix(input, "data: [DONE]")
}
