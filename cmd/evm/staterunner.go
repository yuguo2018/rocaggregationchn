package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/tests"
	"github.com/urfave/cli/v2"
)

var stateTestCommand = &cli.Command{
	Action:    stateTestCmd,
	Name:      "statetest",
	Usage:     "Executes the given state tests. Filenames can be fed via standard input (batch mode) or as an argument (one-off execution).",
	ArgsUsage: "<file>",
}

// StatetestResult contains the execution status after running a state test, any
// error that might have occurred and a dump of the final state if requested.
type StatetestResult struct {
	Name  string       `json:"name"`
	Pass  bool         `json:"pass"`
	Root  *common.Hash `json:"stateRoot,omitempty"`
	Fork  string       `json:"fork"`
	Error string       `json:"error,omitempty"`
	State *state.Dump  `json:"state,omitempty"`
}

func stateTestCmd(ctx *cli.Context) error {
	// Configure the EVM logger
	config := &logger.Config{
		EnableMemory:     !ctx.Bool(DisableMemoryFlag.Name),
		DisableStack:     ctx.Bool(DisableStackFlag.Name),
		DisableStorage:   ctx.Bool(DisableStorageFlag.Name),
		EnableReturnData: !ctx.Bool(DisableReturnDataFlag.Name),
	}
	var cfg vm.Config
	switch {
	case ctx.Bool(MachineFlag.Name):
		cfg.Tracer = logger.NewJSONLogger(config, os.Stderr)

	case ctx.Bool(DebugFlag.Name):
		cfg.Tracer = logger.NewStructLogger(config)
	}
	// Load the test content from the input file
	if len(ctx.Args().First()) != 0 {
		return runStateTest(ctx.Args().First(), cfg, ctx.Bool(MachineFlag.Name), ctx.Bool(DumpFlag.Name))
	}
	// Read filenames from stdin and execute back-to-back
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fname := scanner.Text()
		if len(fname) == 0 {
			return nil
		}
		if err := runStateTest(fname, cfg, ctx.Bool(MachineFlag.Name), ctx.Bool(DumpFlag.Name)); err != nil {
			return err
		}
	}
	return nil
}

// runStateTest loads the state-test given by fname, and executes the test.
func runStateTest(fname string, cfg vm.Config, jsonOut, dump bool) error {
	src, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	var testsByName map[string]tests.StateTest
	if err := json.Unmarshal(src, &testsByName); err != nil {
		return err
	}

	// Iterate over all the tests, run them and aggregate the results
	results := make([]StatetestResult, 0, len(testsByName))
	for key, test := range testsByName {
		for _, st := range test.Subtests() {
			// Run the test and aggregate the result
			result := &StatetestResult{Name: key, Fork: st.Fork, Pass: true}
			test.Run(st, cfg, false, rawdb.HashScheme, func(err error, tstate *tests.StateTestState) {
				var root common.Hash
				if tstate.StateDB != nil {
					root = tstate.StateDB.IntermediateRoot(false)
					result.Root = &root
					if jsonOut {
						fmt.Fprintf(os.Stderr, "{\"stateRoot\": \"%#x\"}\n", root)
					}
					if dump { // Dump any state to aid debugging
						cpy, _ := state.New(root, tstate.StateDB.Database(), nil)
						dump := cpy.RawDump(nil)
						result.State = &dump
					}
				}
				if err != nil {
					// Test failed, mark as so
					result.Pass, result.Error = false, err.Error()
				}
			})
			results = append(results, *result)
		}
	}
	out, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(out))
	return nil
}
