package main

import (
	"github.com/eyecuelab/go-api/cmd"
	// "github.com/eyecuelab/go-api/cmd/env"
)

func main() {
	cmd.Exec()
}

func init() {
	// env.Load(os.Getenv("ENV"))
	cmd.Init()
}
