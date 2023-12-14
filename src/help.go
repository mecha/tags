package tags

import (
	"fmt"
	"os"
)

func Help() string {
    s := fmt.Sprintln("Tags - A rule based directory tagging utility.")
	s += fmt.Sprintln("https://github.com/mecha/tags")
	s += fmt.Sprintln()
	s += fmt.Sprintln("  Usage:")
	s += fmt.Sprintf("    %s [options] [directory]\n", os.Args[0])
	s += fmt.Sprintln()
	s += fmt.Sprintf("  Arguments:\n")
	s += fmt.Sprintf("    directory\t\tThe directory to search for tags. Defaults to the current directory.")
	s += fmt.Sprintln()
	s += fmt.Sprintln("  Options:")
	s += fmt.Sprintf("    -h, --help\t\tShow this help message\n")
	s += fmt.Sprintf("    -c, --config\tThe path to the config file. Default: %s\n", DefaultConfigPath())
	s += fmt.Sprintf("    -v, --verbose\tShow verbose output\n")

    return s
}
