package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/urfave/cli/v2"
)

const (
	defaultKeyfileName = "keyfile.json"
)

var app *cli.App

func init() {
	app = flags.NewApp("Ethereum key manager")
	app.Commands = []*cli.Command{
		commandGenerate,
		commandInspect,
		commandChangePassphrase,
		commandSignMessage,
		commandVerifyMessage,
	}
}

// Commonly used command line flags.
var (
	passphraseFlag = &cli.StringFlag{
		Name:  "passwordfile",
		Usage: "the file that contains the password for the keyfile",
	}
	jsonFlag = &cli.BoolFlag{
		Name:  "json",
		Usage: "output JSON instead of human-readable format",
	}
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
