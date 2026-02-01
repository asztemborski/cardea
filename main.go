package main

import (
	"context"
	"os"

	"github.com/asztemborski/cardea/cmd"
)

func main() {
	cmd.Execute(context.Background(), os.Args)
}
