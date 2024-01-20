package main

import (
	"os"

	"github.com/curtisnewbie/hammer/internal/hammer"
)

func main() {
	hammer.BootstrapServer(os.Args)
}
