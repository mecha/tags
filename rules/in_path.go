package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mecha/tags/log"
	"github.com/mecha/tags/utils"
)

type (
	InPath struct {
		Paths []string
	}
)

func (r *InPath) Load(cfg map[string]interface{}) error {
	if r.Paths == nil {
		r.Paths = make([]string, 0)
	}

	pathsVal, ok := cfg["paths"]
	if !ok {
		return fmt.Errorf("[in_path] no \"paths\" key in config\n")
	}

	list, ok := pathsVal.([]interface{})
	if !ok {
		return fmt.Errorf("[in_path] \"paths\" is not a list: %v", list)
	}

	for _, path := range list {
		path, ok := path.(string)

		if !ok {
			return fmt.Errorf("[in_path] invalid value: %v\n", path)
		}

		r.Paths = append(r.Paths, path)
	}

	return nil
}

func (r *InPath) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"paths": r.Paths,
	}
}

func (r *InPath) Evaluate(dir string) (bool, error) {
	for _, path := range r.Paths {
		log.Debug("   [in_path] %s\n", path)

		ePath, err := utils.ExpandTilde(path)

		if err != nil {
			return false, err
		}

		abs, err := filepath.Abs(ePath)

		if err != nil {
			return false, err
		} else if strings.HasPrefix(dir, abs) {
			return true, nil
		}
	}

	return false, nil
}

func (r *InPath) Add(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No paths specified.")
	}

	for _, arg := range args {
		r.Paths = append(r.Paths, arg)
	}

	return nil
}

func (r *InPath) Del(args []string) error {
	if len(args) == 0 {
		r.Paths = make([]string, 0)
		return nil
	}

	for _, arg := range args {
		for i, path := range r.Paths {
			if path == arg {
				last := len(r.Paths) - 1
				r.Paths[i] = r.Paths[last]
				r.Paths = r.Paths[:last]
			}
		}
	}

	return nil
}

func (r *InPath) String() string {
	s := ""
	for _, path := range r.Paths {
		s += fmt.Sprintf("\n[in_path] %s", path)
	}

	return strings.TrimLeft(s, "\n")
}
