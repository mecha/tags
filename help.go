package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mecha/tags/config"
)

func helpCommand(args []string) {
	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "find":
		printFindHelp()
	case "show":
		printShowHelp()
	case "add":
		printAddHelp()
	case "rm":
		printRmHelp()
	case "rules":
		printRulesHelp()
	case "config":
		printConfigHelp()
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Printf(`Tags - Use rules to tag directories.
https://github.com/mecha/tags

USAGE:
  %s [<COMMAND>] [<ARGUMENTS>] [<OPTIONS>]

COMMANDS
  find          (Default) Output the tags for the given/current directory.
  show          Show all the tags and their rules. 
  add           Add new tags or rules.
  rm            Remove a tag or rule.
  help          Show this help message.

OPTIONS
`, os.Args[0])

	printOptions()

	fmt.Printf(`
HELP PAGES

  View these help pages using the "%[1]s help <PAGE>" command.

  <COMMAND>   Help about a specific command.
  rules       Rules and the available rule types.
  config      Information about the config file.
`, os.Args[0])
}

func printOptions() {
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  -%s\t\t%s\n", f.Name, f.Usage)
	})
}

func printFindHelp() {
	fmt.Printf(`DESCRIPTION

  Output the tags that match the given/current directory.
  This is the default command. The "find" may be omitted.

SYNOPSIS

  %[1]s [find] [<DIRECTORY>] [<OPTIONS>]

ARGUMENTS

  <DIRECTORY>   The directory to find tags for. Defaults to the current directory.
                Default: %[2]s

OPTIONS

`, os.Args[0], configPath)

	printOptions()

	fmt.Printf(`
EXAMPLES

  %[1]s
  %[1]s ~/Documents
  %[1]s my/project -p
`, os.Args[0])
}

func printShowHelp() {
	fmt.Printf(`DESCRIPTION

  Shows all the tags and their rules.

SYNOPSIS

  %[1]s show [<OPTIONS>]

OPTIONS
`, os.Args[0], configPath)

	printOptions()

	fmt.Printf(`
EXAMPLES

  %[1]s show
  %[1]s show -c ~/backup-rules.json
`, os.Args[0])
}

func printAddHelp() {
	fmt.Printf(`DESCRIPTION

  Adds a new tag rule.

SYNOPSIS

  %[1]s add <TAG> <RULETYPE> <VALUES>... [<OPTIONS>]

ARGUMENTS

  <TAG>         The tag to add the rule to.
  <RULETYPE>    The type of rule to add. See "%[1]s help rules" for more info.
  <VALUES>...   The values to add to the rule.

OPTIONS

`, os.Args[0], configPath)

	printOptions()

	fmt.Printf(`
EXAMPLES

  %[1]s add make file_exists Makefile
  %[1]s add react file_contains pakage.json react
  %[1]s add tmp in_path /tmp
`, os.Args[0])
}

func printRmHelp() {
	fmt.Printf(`DESCRIPTION

  Removes an entire tag, or rules from a tag.

SYNOPSIS

  %[1]s rm <TAG> [<RULETYPE>] [<VALUES>...] [<OPTIONS>]

ARGUMENTS

  <TAG>         The tag to remove, or remove rules from.
  <RULETYPE>    The type of rule to remove, or remove from.
                If omitted, the entire tag will be removed.
                See "%[1]s help rules" for more info.
  <VALUES>...   The values to add to the rule.
                If omitted, all rules of the given <RULETYPE> will be removed.

OPTIONS

`, os.Args[0], configPath)

	printOptions()

	fmt.Printf(`
EXAMPLES

  Remove an entire tag:
    %[1]s rm make

  Remove all rules of a specific type:
    %[1]s rm go file_exists

  Remove a specific rule:
    %[1]s rm go file_exists go.mod
`, os.Args[0])
}

func printRulesHelp() {
	fmt.Printf(`RULE TYPES

This help page details the available tag rule types, how they match directories,
and how to configure them.

================================================================================
file_exists

  Rules with this type match directories that have a specific file.
  The file can be either a file or a directory, and can be deeply nested
  within the directory. The path is relative to the directory being checked.

  Rules of this type take only one argument: the path to the file.

  Add:        %[1]s add <TAG> file_exists <FILE>
  Remove:     %[1]s rm <TAG> file_exists <FILE>

  Examples:   %[1]s add make file_exists Makefile
              %[1]s rm make file_exists Makefile

  The files that the rule matches with are stored in a "files" list in the
  config. Example:

  ┌─ rules.json ──────────────────────┐
  │ {                                 │
  │   "make": {                       │
  │     "file_exists": {              │
  │       "files": [                  │
  │         "Makefile"                │
  │       ]                           │
  │     }                             │
  │   }                               │
  │ }                                 │
  └───────────────────────────────────┘ 

================================================================================
file_contains

  Rules with this type match directories that have a specific file, which also
  contains a specific substring. The search is performed on the entire file,
  using a case-sensitive, non-whole-word search.

  Rules of this type take two arguments: the path to the file and the substring
  to search for.

  Add:        %[1]s add <TAG> file_contains <FILE> <TEXT>
  Remove:     %[1]s rm <TAG> file_contains <FILE> <TEXT>

  Examples:   %[1]s add react file_contains package.json react
              %[1]s rm react file_contains package.json react

  This rule type stores the files and the substrings as a "search" object in
  the config file. The keys of the object are the file names, and the values
  are the search substrings. Example:

  ┌─ rules.json ──────────────────────┐
  │ {                                 │
  │   "react": {                      │
  │     "file_contains": {            │
  │       "search": {                 │
  │         "package.json": "react"   │
  │       }                           │
  │     }                             │
  │   }                               │
  │ }                                 │
  └───────────────────────────────────┘ 

================================================================================
in_path

  Rules with this type match directories that are in a specific path. The path
  is checked using a prefix match, which means that the prefix path must be
  absolute (either relative to the root or some other path token, such as "~"
  or "%%APPDATA%%").

  Rules of this type take only one argument: the path to the directory.

  Add:        %[1]s add <TAG> in_path <PATH>
  Remove:     %[1]s rm <TAG> in_path <PATH>

  Examples:   %[1]s add fonts in_path /usr/share/fonts
              %[1]s rm react file_contains package.json react

  In the config, these rules store the paths in a "paths" list. Example:

  ┌─ rules.json ──────────────────────┐
  │ {                                 │
  │   "fonts": {                      │
  │     "in_path": {                  │
  │       "paths": [                  │
  │         "/usr/share/fonts"        │
  │       ]                           │
  │     }                             │
  │   }                               │
  │ }                                 │
  └───────────────────────────────────┘ 

================================================================================
MORE HELP

  %[1]s help add
  %[1]s help rm
  %[1]s help config
`, os.Args[0])
}

func printConfigHelp() {
	fmt.Printf(`THE CONFIG FILE

  Default: %s
  
  The config file stores all of your tags and their rules in JSON format, using
  the below structure:
  
  ┌────────────────────────────┐
  │ {                          │
  │   "tag1": {                │
  │     "rule_type": {         │
  │       <rule type config>   │
  │     }                      │
  │   }                        │
  │ }                          │
  └────────────────────────────┘
  
  The root object contains the tags. The name of each property in the object is
  the name of the tag, and the value is the tag's rules object.
  
  A rules object consists of one property for each rule type. Not all rule types
  need to be included. Each property's name is the rule type, and the value is
  the rule type's config object. It's structure depends on the rule type.
  
  See "%s help rules" for more information about the rule types.

MULTIPLE CONFIG FILES

  When running %[2]s, your main config file will be used by default. You can
  specify a different config file using the "-c" option. For example:
  
    %[2]s -c ~/backup-rules.json

  Alternatively, you can set the "TAGS_CONFIG" environment variable to the path
  of the config file you want to use. For example:

    TAGS_CONFIG="~/backup-rules.json" %[2]s

`, config.DefaultPath(), os.Args[0])
}
