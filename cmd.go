package main

import "flag"

type Cmd struct {
	dbPath       string
	staticPath   string
	addr         string
	ipHttpHeader string
	debug        bool
}

func parseCmd() Cmd {
	var cmd Cmd
	flag.StringVar(&cmd.dbPath, "data.path", "./notepad.db", "database save dbPath")
	flag.StringVar(&cmd.staticPath, "static.path", "./", "static path")
	flag.StringVar(&cmd.addr, "server.addr", "0.0.0.0:8083", "addr")
	flag.StringVar(&cmd.ipHttpHeader, "http.header.ip", "NOT", "[NOT|X-Forwarded-For|X-Real-Ip|X-Real-Forwarded-For|X-Real-Forwarded-For...]")
	flag.BoolVar(&cmd.debug, "debug", false, "debug")
	flag.Parse()
	return cmd
}
