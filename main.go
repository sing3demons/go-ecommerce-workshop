package main

import (
	"fmt"
	"os"

	"github.com/sing3demons/shop/config"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env.dev"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())
	fmt.Println(cfg.App().Name())
}
