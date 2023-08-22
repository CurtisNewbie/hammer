package main

import (
	"os"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/curtisnewbie/hammer/hammer"
)

func main() {
	c := common.EmptyRail()
	hammer.PrepareServer(c)
	server.BootstrapServer(os.Args)
}
