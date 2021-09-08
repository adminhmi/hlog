package cmd

import "golang.org/x/sys/unix"

const ioctlReadTerms = unix.TCGETS

func isTerminal(fd int) bool {
	_, err := unix.IoctlGetTermios(fd, ioctlReadTerms)
	return err == nil
}
