package main

import (
	"github.com/valyala/fasthttp"
	"github.com/asdine/storm"
)

var (
	cmd Cmd
	db  *storm.DB
)

func main() {
	cmd = parseCmd()
	RestoreAssets(cmd.staticPath, "static")
	db, _ = storm.Open(cmd.dbPath)
	defer db.Close()
	fasthttp.ListenAndServe(cmd.addr, fastHTTPHandler)
}
