package logging

import (
	"github.com/phuslu/log"
)

func Setup(level string, format string) {
	logLevel := log.ParseLevel(level)
	if logLevel == 8 {
		log.Warn().Msg("Invalid log level, defaulting to info")
		logLevel = log.InfoLevel
	}

	// Set log format
	if format == "json" {
		log.DefaultLogger = log.Logger{
			Level:      logLevel,
			TimeFormat: "2006-01-02 15:04:05",
			Caller:     1,
		}
	} else {
		// Default to text format
		log.DefaultLogger = log.Logger{
			TimeFormat: "15:04:05",
			Caller:     1,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}

}
