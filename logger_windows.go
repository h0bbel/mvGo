//go:build windows
// +build windows

package main

import (
	"io"
	"log"
)

func newLogger(outputs []io.Writer, sysCfg SyslogConfig) *log.Logger {
	// Syslog not supported on Windows
	return log.New(io.MultiWriter(outputs...), "", log.LstdFlags)
}
