package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/mecha/tags/config"
	"github.com/mecha/tags/log"
	tags "github.com/mecha/tags/tags"
)

var (
	configPath string
	directory  string
	help       bool
	parallel   bool
	quiet      bool
	verbose    bool
	verbose2   bool
)

func main() {
	flag.Usage = printHelp

	flag.StringVar(&configPath, "c", config.DefaultPath(), "The path to the config file.")
	flag.BoolVar(&parallel, "p", false, "Execute tag rules in parallel.")
	flag.BoolVar(&quiet, "q", false, "Suppress all output.")
	flag.BoolVar(&verbose, "v", false, "Show verbose output.")
	flag.BoolVar(&verbose2, "vv", false, "Show debugging output.")
	flag.Parse()

	switch true {
	case verbose2:
		log.SetLevel(log.DebugLevel)
	case verbose:
		log.SetLevel(log.InfoLevel)
	case quiet:
		log.SetLevel(log.QuietLevel)
	}

	log.Info("Reading config from %s\n", configPath)

	cfg, err := config.Read(configPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file. %s\n", err)
		os.Exit(1)
	}

	log.Debug("Read %d tags\n", len(cfg))

	command := ""
	args := flag.Args()

	if len(args) > 0 {
		command = args[0]
	}

	switch command {
	case "help":
        helpCommand(args[1:])
	case "add":
		log.Debug("Running `add` command\n")
        addCommand(cfg, args[1:])
	case "rm":
		log.Debug("Running `rm` command\n")
        rmCommand(cfg, args[1:])
	case "show":
		log.Debug("Running `show` command\n")
		showCommand(cfg)
	case "find":
        args = args[1:]
        fallthrough
    default:
		log.Debug("Running `find` command\n")
		findCommand(cfg, args)
	}

	os.Exit(0)
}

func findCommand(cfg map[string]tags.Tag, args []string) {
	dir := ""

	if len(args) == 0 {
		log.Debug("No directory specified, using current directory\n")
		if cwd, err := os.Getwd(); err == nil {
			dir = cwd
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		dir = args[0]
	}

	log.Info("Directory = %s\n", dir)

	if parallel {
		log.Info("Checking tag rules in parallel\n")
		wg := sync.WaitGroup{}

		for tagName, tag := range cfg {
			wg.Add(1)

			go func(name string, tag tags.Tag) {
				defer wg.Done()
				log.Debug("=> %s\n", name)

				if tags.IsMatch(dir, &tag) {
					fmt.Println(name)
				}
			}(tagName, tag)
		}

		wg.Wait()
	} else {
		log.Info("Checking tag rules in series\n")

		for tagName, tag := range cfg {
			log.Debug("=> %s\n", tagName)

			if tags.IsMatch(dir, &tag) {
				fmt.Println(tagName)
			}
		}

		log.Debug("Done\n")
	}

	os.Exit(0)
}

func addCommand(cfg map[string]tags.Tag, args []string) {
	nArgs := len(args)

	switch true {
	case nArgs == 0:
		log.Error("No tag specified.\n")
		os.Exit(1)
	case nArgs == 1:
		log.Error("No rule type specified.\n")
		os.Exit(1)
	case nArgs == 2:
		log.Error("No rule values specified.\n")
		os.Exit(1)
	}

	tagName := args[0]
	ruleType := args[1]
	values := args[2:]

	log.Debug("Looking for tag with name \"%s\".\n", tagName)

	tag, hasTag := cfg[tagName]
	if !hasTag {
		log.Debug("Not found. Creating new tag.\n")
		tag = tags.Tag{}
	}

	err := tag.AddRule(ruleType, values)

	if err != nil {
		log.Error("%s\n", err)
		os.Exit(1)
	}

	cfg[tagName] = tag
	err = config.Write(configPath, cfg)

	if err != nil {
		log.Error("%s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func showCommand(cfg map[string]tags.Tag) {
	for tagName, tag := range cfg {
		fmt.Printf("[%s]\n", tagName)

		if len(tag.Rules) > 0 {
			for _, rule := range tag.Rules {
				fmt.Println(rule)
			}

			fmt.Println()
		}
	}
}

func rmCommand(cfg map[string]tags.Tag, args []string) {
	nArgs := len(args)

	switch true {
	case nArgs == 0:
		log.Error("No tag specified.\n")
		os.Exit(1)

	case nArgs == 1:
		tagName := args[0]
		log.Info("Removing tag \"%s\"\n", tagName)
		delete(cfg, tagName)

	default:
		tagName := args[0]
		ruleType := args[1]

		tag, ok := cfg[tagName]
		if !ok {
			log.Error("Tag \"%s\" not found\n", tagName)
			os.Exit(1)
		}

		numDel := nArgs - 2
		if numDel <= 0 {
			numDel = len(cfg[tagName].Rules)
		}

		log.Info("Removing %d \"%s\" rules from the \"%s\" tag\n", numDel, ruleType, tagName)

		tag.DelRule(ruleType, nil)
		cfg[tagName] = tag
	}

	err := config.Write(configPath, cfg)

	if err != nil {
		log.Error("%s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
