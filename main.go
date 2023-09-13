package main

import (
	"os"

	"github.com/curtisnewbie/hammer/hammer"
	"github.com/curtisnewbie/miso/miso"
)

func main() {
	miso.PreServerBootstrap(func(rail miso.Rail) error {
		hammer.PrepareServer(rail)
		return nil
	})

	miso.BootstrapServer(os.Args)
}
