package main

import (
	"fmt"
	"os"
	"strings"
)

// TODO - UGLY
func printUserdocs(filename string, userdocs []userdoc) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.WriteString("| Output | Definition | Selector | Location |\n| ---- | ---- | ---- |\n")
	if err != nil {
		return err
	}

	for _, u := range userdocs {
		_, err = f.WriteString(fmt.Sprintf("| %s | %s | %s %s |\n", strings.Join(u.Params, ", "), strings.Join(u.Comments, ", "), u.selector, u.lineNumber))
		if err != nil {
			return err
		}
	}
	return nil
}
