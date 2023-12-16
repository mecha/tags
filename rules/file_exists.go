package rules

import (
	"fmt"
	"os"
	"strings"

	"github.com/mecha/tags/log"
)

type (
	FileExists struct {
		Files []string
	}
)

func (rule *FileExists) Load(cfg map[string]interface{}) error {
	if rule.Files == nil {
		rule.Files = make([]string, 0)
	}

	filesVal, ok := cfg["files"]
	if !ok {
		return fmt.Errorf("[file_exists] no \"files\" key in config\n")
	}

	list, ok := filesVal.([]interface{})
	if !ok {
		return fmt.Errorf("[file_exists] \"files\" is not a list: %v", filesVal)
	}

	for _, file := range list {
		file, ok := file.(string)

		if !ok {
			return fmt.Errorf("[file_exists] invalid value: %v\n", file)
		}

		rule.Files = append(rule.Files, file)
	}

	return nil
}

func (rule *FileExists) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"files": rule.Files,
	}
}

func (rule *FileExists) Evaluate(dir string) (bool, error) {
	for _, file := range rule.Files {
		log.Debug("  file_exists: %s\n", file)

		fullpath := dir + "/" + file
		_, err := os.Stat(fullpath)

		if err == nil {
			return true, nil
		} else if os.IsNotExist(err) {
			continue
		} else {
			return false, err
		}
	}

	return false, nil
}

// tags add rust file_exists Cargo.toml
func (rule *FileExists) Add(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No files specified.")
	}

	for _, arg := range args {
		rule.Files = append(rule.Files, arg)
	}

	return nil
}

func (rule *FileExists) Del(args []string) error {
	if len(args) == 0 {
		rule.Files = make([]string, 0)
		return nil
	}

	for _, arg := range args {
		for i, file := range rule.Files {
			if file == arg {
				last := len(rule.Files) - 1
				rule.Files[i] = rule.Files[last]
				rule.Files = rule.Files[:last]
			}
		}
	}

	return nil
}

func (r *FileExists) String() string {
	s := ""
	for _, file := range r.Files {
		s += fmt.Sprintf("\nfile_exists: %s", file)
	}

	return strings.TrimLeft(s, "\n")
}
