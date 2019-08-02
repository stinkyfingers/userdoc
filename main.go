package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	files     = flag.String("f", "", "files (or dirs), comma-separated")
	output    = flag.String("o", "output.md", "output file")
	selectors = flag.String("s", "", "additional selectors, comma-separated (e.g. 'foo.Bar')")
)

func main() {
	flag.Parse()
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var userdocs []userdoc

	for _, file := range strings.Split(*files, ",") {
		err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) != ".go" {
				return nil
			}

			userdoc, err := parseFile(path, strings.Split(*selectors, ","))
			if err != nil {
				return err
			}
			userdocs = append(userdocs, userdoc...)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return printUserdocs(*output, userdocs)
}
