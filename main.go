package main

import (
	"os"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/curtisnewbie/hammer/hammer"
)

func main() {
	c := common.EmptyExecContext()
	hammer.PrepareServer(c)
	server.DefaultBootstrapServer(os.Args, c)
}
