package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type InPathRule struct {
	paths []string
}

func (rule *InPathRule) SetConfig(config RuleConfig) error {
	rule.paths = make([]string, 0)

	switch config.(type) {
	case string:
		rule.paths = append(rule.paths, config.(string))

	case []interface{}:
		for _, file := range config.([]interface{}) {
			rule.paths = append(rule.paths, file.(string))
		}

	default:
		return fmt.Errorf("Config for \"files_exist\" rule must be a string or list of strings, %T given", config)
	}

	return nil
}

func (rule *InPathRule) Evaluate(dir string, app App) bool {
	for _, path := range rule.paths {
        path = expandTilde(path)
		abs, err := filepath.Abs(path)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		} else if strings.HasPrefix(dir, abs) {
			return true
		}

	}

	return false
}
