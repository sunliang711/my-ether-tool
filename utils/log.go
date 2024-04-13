package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var Logger zerolog.Logger

func SetLogger(cmd *cobra.Command) error {
	loglevel, err := zerolog.ParseLevel(cmd.Flag("loglevel").Value.String())
	if err != nil {
		return fmt.Errorf("parse loglevel error: %w", err)
	}
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s |", i)
		},
		FormatCaller: func(i interface{}) string {
			paths := strings.Split(i.(string), "/")
			l := len(paths)
			if l > 2 {
				return strings.Join([]string{paths[l-2], paths[l-1]}, "/")
			}
			return filepath.Base(fmt.Sprintf("%s", i))
		},
		PartsExclude: []string{
			// zerolog.TimestampFieldName,
		},
	}

	Logger = zerolog.New(writer).Level(loglevel).With().Timestamp().Caller().Logger()

	return nil
}

func GetLogger(funcName string) zerolog.Logger {
	return Logger.With().Str("function", funcName).Logger()
}
