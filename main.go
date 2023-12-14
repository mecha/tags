package main

import (
	"flag"
	"fmt"
	"os"

	tags "github.com/mecha/tags/src"
)

var (
	help       bool
	verbose    bool
	configPath string
	directory  string
)

func main() {
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.StringVar(&configPath, "c", tags.DefaultConfigPath(), "Config path")

	flag.Parse()

	if help {
		fmt.Println(tags.Help())
		os.Exit(0)
	}

	directory := flag.Arg(0)
	if directory == "" {
		if cwd, err := os.Getwd(); err == nil {
			directory = cwd
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	t, err := tags.New(configPath, verbose)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	t.FindTags(directory, func(tag tags.Tag) {
		fmt.Println(tag.Name)
	})

	os.Exit(0)
}
