//go:build !windows
// +build !windows

package main

import (
	"io"
	"log"
	"log/syslog"
)

func newLogger(outputs []io.Writer, sysCfg SyslogConfig) *log.Logger {
	if sysCfg.Enabled && sysCfg.Address != "" {
		network := sysCfg.Network
		if network == "" {
			network = "udp"
		}
		sysLogger, err := syslog.Dial(network, sysCfg.Address, syslog.LOG_INFO|syslog.LOG_USER, "mvGo")
		if err == nil {
			outputs = append(outputs, sysLogger)
		} else {
			log.Println("Warning: failed to connect to syslog:", err)
		}
	}
	return log.New(io.MultiWriter(outputs...), "", log.LstdFlags)
}
