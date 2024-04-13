package utils

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func ExitWithMsgWhen(condition bool, format string, args ...any) {
	if condition {
		fmt.Fprintf(os.Stderr, format, args...)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
}

func ExitWhenError(err error, format string, args ...any) {
	ExitWithMsgWhen(err != nil, format, args...)
}

func ExitWhen(logger zerolog.Logger, condition bool, format string, args ...any) {
	if condition {
		logger.Fatal().Msgf(format, args...)
	}
}

func ExitWhenErr(logger zerolog.Logger, err error, format string, args ...any) {
	ExitWhen(logger, err != nil, format, args...)
}
