package utils

import (
	"fmt"
	"os"
)

func ExitWithMsgWhen(condition bool, format string, args ...any) {
	if condition {
		fmt.Printf(format, args...)
		os.Exit(1)
	}
}

func ExitWhenError(err error, format string, args ...any) {
	ExitWithMsgWhen(err != nil, format, args...)
}
