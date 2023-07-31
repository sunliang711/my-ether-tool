package utils

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func ReadSecret(prompt string) (string, error) {
	fmt.Fprintf(os.Stderr, "%s", prompt)
	secret, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Fprint(os.Stderr, "\n")
	return string(secret), nil
}

func ReadChar(prompt string) (byte, error) {
	fmt.Fprintf(os.Stderr, "%s", prompt)
	rd := bufio.NewReader(os.Stdin)
	return rd.ReadByte()
}
