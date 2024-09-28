package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/core/asm"
	"github.com/urfave/cli/v2"
)

var disasmCommand = &cli.Command{
	Action:    disasmCmd,
	Name:      "disasm",
	Usage:     "Disassembles evm binary",
	ArgsUsage: "<file>",
}

func disasmCmd(ctx *cli.Context) error {
	var in string
	switch {
	case len(ctx.Args().First()) > 0:
		fn := ctx.Args().First()
		input, err := os.ReadFile(fn)
		if err != nil {
			return err
		}
		in = string(input)
	case ctx.IsSet(InputFlag.Name):
		in = ctx.String(InputFlag.Name)
	default:
		return errors.New("missing filename or --input value")
	}

	code := strings.TrimSpace(in)
	fmt.Printf("%v\n", code)
	return asm.PrintDisassembled(code)
}
