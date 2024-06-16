// Project entry program.
package main

import (
	"slender/internal/daemon"
	"slender/internal/parse"
)

func main() {
	parse.Parse()
	daemon.StartDaemon()
}
