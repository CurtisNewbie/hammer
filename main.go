package main

import (
	"os"

	"github.com/curtisnewbie/hammer/hammer"
	"github.com/curtisnewbie/miso/core"
	"github.com/curtisnewbie/miso/server"
)

func main() {
	server.PreServerBootstrap(func(rail core.Rail) error {
		hammer.PrepareServer(rail)
		return nil
	})

	server.BootstrapServer(os.Args)
}
