package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("usage: %s <config.toml>", os.Args[0]))
	}

	cfgFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer cfgFile.Close()

	cfg, err := parseConfig(cfgFile)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	run(ctx, cfg)
}
