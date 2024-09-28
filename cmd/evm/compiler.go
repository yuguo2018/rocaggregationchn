package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/cmd/evm/internal/compiler"

	"github.com/urfave/cli/v2"
)

var compileCommand = &cli.Command{
	Action:    compileCmd,
	Name:      "compile",
	Usage:     "Compiles easm source to evm binary",
	ArgsUsage: "<file>",
}

func compileCmd(ctx *cli.Context) error {
	debug := ctx.Bool(DebugFlag.Name)

	if len(ctx.Args().First()) == 0 {
		return errors.New("filename required")
	}

	fn := ctx.Args().First()
	src, err := os.ReadFile(fn)
	if err != nil {
		return err
	}

	bin, err := compiler.Compile(fn, src, debug)
	if err != nil {
		return err
	}
	fmt.Println(bin)
	return nil
}
