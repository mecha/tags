package tags

import (
	"fmt"
	"os"
)

type FileExistsRule struct {
	files []string
}

func (rule *FileExistsRule) SetConfig(config RuleConfig) error {
	rule.files = make([]string, 0)

	switch config.(type) {
	case string:
		rule.files = append(rule.files, config.(string))

	case []interface{}:
		for _, file := range config.([]interface{}) {
			rule.files = append(rule.files, file.(string))
		}

	default:
		return fmt.Errorf("Config for \"files_exist\" rule must be a string or list of strings, %T given", config)
	}

	return nil
}

func (rule FileExistsRule) Evaluate(dir string, app App) bool {
	for _, file := range rule.files {
		app.log("=> [file_exists] %s ... ", file)

		fullpath := dir + "/" + file
		_, err := os.Stat(fullpath)

		if err == nil {
			app.log("yes\n")
			return true
		} else if os.IsNotExist(err) {
			app.log("no\n")
			continue
		} else {
			app.log("\n")
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}

	return false
}
