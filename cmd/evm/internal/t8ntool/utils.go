package t8ntool

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// readFile reads the json-data in the provided path and marshals into dest.
func readFile(path, desc string, dest interface{}) error {
	inFile, err := os.Open(path)
	if err != nil {
		return NewError(ErrorIO, fmt.Errorf("failed reading %s file: %v", desc, err))
	}
	defer inFile.Close()
	decoder := json.NewDecoder(inFile)
	if err := decoder.Decode(dest); err != nil {
		return NewError(ErrorJson, fmt.Errorf("failed unmarshaling %s file: %v", desc, err))
	}
	return nil
}

// createBasedir makes sure the basedir exists, if user specified one.
func createBasedir(ctx *cli.Context) (string, error) {
	baseDir := ""
	if ctx.IsSet(OutputBasedir.Name) {
		if base := ctx.String(OutputBasedir.Name); len(base) > 0 {
			err := os.MkdirAll(base, 0755) // //rw-r--r--
			if err != nil {
				return "", err
			}
			baseDir = base
		}
	}
	return baseDir, nil
}
