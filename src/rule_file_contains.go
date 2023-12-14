package tags

import (
	"fmt"
	"os"
	"strings"
)

type FileContainsRule struct {
	Mappings map[string]string
}

func (rule *FileContainsRule) SetConfig(config RuleConfig) error {
	if config, ok := config.(map[string]interface{}); ok {
		rule.Mappings = make(map[string]string)

		for key, value := range config {
			rule.Mappings[key] = value.(string)
		}
	} else {
		return fmt.Errorf("Config for \"file_contains\" rule must be an object, %T given.", config)
	}

	return nil
}

func (rule FileContainsRule) Evaluate(dir string, app App) bool {
	for file, text := range rule.Mappings {
		app.log("=> [file_contains] %s contains `%s` ... ", file, text)

		contents, err := os.ReadFile(dir + "/" + file)

		if err == nil {
			if strings.Contains(string(contents), text) {
				app.log("yes\n")
				return true
			} else {
				app.log("no\n")
			}
		} else {
			app.log("no\n")
		}
	}

	return false
}
