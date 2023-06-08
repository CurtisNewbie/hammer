package main

import (
	"os"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
)

func main() {
	c := common.EmptyExecContext()
	server.DefaultBootstrapServer(os.Args, c)
}
