package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("usage: %s <config.toml>", os.Args[0]))
	}

	configFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	var cfg Config
	if err := toml.NewDecoder(configFile).Decode(&cfg); err != nil {
		panic(err)
	}

	ctx := context.Background()

	run(ctx, cfg)
}
