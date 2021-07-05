package hlog

import (
	"io"
)

func CheckIfTerminal(w io.Writer) bool {
	return true
}
