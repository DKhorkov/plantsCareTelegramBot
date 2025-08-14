package utils

import "fmt"

const hashLength = 8

// GenUniqueParam Generates a InlineButton.Unique param via concatenating
// base name with a random string.
// It is necessary to exclude the overlap of links between buttons.
func GenUniqueParam(base string) string {
	return fmt.Sprintf("%s_%s", base, RandSequence(hashLength))
}
