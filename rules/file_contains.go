package rules

import (
	"fmt"
	"os"
	"strings"

	"github.com/mecha/tags/log"
)

type (
	FileContains struct {
		Search map[string]string
	}
)

func (r *FileContains) Load(cfg map[string]interface{}) error {
	if r.Search == nil {
		r.Search = make(map[string]string)
	}

	searchVal, ok := cfg["search"]
	if !ok {
		return fmt.Errorf("[file_contains] missing \"search\" in config\n")
	}

	dict, ok := searchVal.(map[string]interface{})
	if !ok {
		return fmt.Errorf("[file_contains] \"search\" is not an object: %v", searchVal)
	}

	for file, text := range dict {
		text, ok := text.(string)

		if !ok {
			return fmt.Errorf("[file_contains] invalid value for file \"%s\": %v\n", file, text)
		}

		r.Search[file] = text
	}

	return nil
}

func (r *FileContains) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"search": r.Search,
	}
}

func (r *FileContains) Evaluate(dir string) (bool, error) {
	for file, text := range r.Search {
		log.Debug("   file_contains: %s >> `%s`\n", file, text)

		contents, err := os.ReadFile(dir + "/" + file)

		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				return false, err
			}
		}

		if strings.Contains(string(contents), text) {
			return true, nil
		}
	}

	return false, nil
}

func (r *FileContains) Add(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No mappings provided")
	} else if len(args)%2 != 0 {
		return fmt.Errorf("Odd number of arguments")
	} else {
		for i := 0; i < len(args); i += 2 {
			r.Search[args[i]] = args[i+1]
		}
	}

	return nil
}

func (r *FileContains) Del(args []string) error {
	nArgs := len(args)

	if nArgs == 0 {
		r.Search = make(map[string]string)
		return nil
	}

	if nArgs%2 != 0 {
		return fmt.Errorf("Odd number of arguments")
	}

	for i := 0; i < nArgs; i += 2 {
		file := args[i]
		text := args[i+1]
		for k, v := range r.Search {

			if k == file && v == text {
				delete(r.Search, k)
			}
		}
	}

	return nil
}

func (r *FileContains) String() string {
	s := ""
	for file, text := range r.Search {
		s += fmt.Sprintf("\nfile_contains: %s >> %s", file, text)
	}

	return strings.TrimLeft(s, "\n")
}
