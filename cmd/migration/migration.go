package main

import (
	"github.com/Toflex/directory_v2/cmd/http/runtime"
	"github.com/Toflex/directory_v2/cmd/migration/atlas"
)

func main() {
	// register providers
	runtime.Register()

	// run ATLAS migration
	atlas.AtlasMigration()
}
