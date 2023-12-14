package tags

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/OpenPeeDeeP/xdg"
)

type (
	App struct {
		Config     []Tag
		ConfigPath string
		Verbose    bool
	}

	Tag struct {
		Name  string
		Rules []TagRule
	}

	TagRule interface {
		SetConfig(RuleConfig) error
		Evaluate(string, App) bool
	}

	ConfigShape map[string]TagConfig
	TagConfig map[string]RuleConfig
	RuleConfig interface{}
)

func New(configPath string, verbose bool) (*App, error) {
	config, err := LoadConfig(configPath)

	if err == nil {
		app := &App{
			Config:     config,
			ConfigPath: configPath,
			Verbose:    verbose,
		}

		app.log("Loaded config from %s (%d tags)\n", configPath, len(config))

		return app, nil
	} else {
		return nil, err
	}
}

func (app App) log(format string, args ...interface{}) {
	if app.Verbose {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func (app App) FindTags(dir string, f func(Tag)) {
	wg := &sync.WaitGroup{}

	for _, tag := range app.Config {
		wg.Add(1)

		if app.Verbose {
			app.CheckTagCallback(dir, tag, wg, f)
		} else {
			go app.CheckTagCallback(dir, tag, wg, f)
		}
	}

	wg.Wait()
}

func (app App) CheckTagCallback(dir string, tag Tag, wg *sync.WaitGroup, f func(Tag)) {
	if wg != nil {
		defer wg.Done()
	}

	if app.CheckTag(dir, tag) {
		f(tag)
	}
}

func (app App) CheckTag(dir string, tag Tag) bool {
	app.log("\nChecking %s tag\n", tag.Name)

	for _, rule := range tag.Rules {
		if rule.Evaluate(dir, app) {
			return true
		}
	}

	return false
}

func DefaultConfigPath() string {
	return xdg.ConfigHome() + "/tags/rules.json"
}

func LoadConfig(path string) ([]Tag, error) {
	_, err := os.Stat(path)
	if err == nil {
		return parseJsonConfig(path)
	} else {
		return make([]Tag, 0), err
	}
}

func parseJsonConfig(path string) ([]Tag, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result ConfigShape
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}

	config := make([]Tag, 0)

	for name, ruleConfigs := range result {
		rules := make([]TagRule, 0)

		for ruleType, ruleConfig := range ruleConfigs {
			rule, err := createRule(ruleType, ruleConfig)
			if err == nil {
				rules = append(rules, rule)
			} else {
				return nil, fmt.Errorf("Tag %s: %s", name, err)
			}
		}

		config = append(config, Tag{name, rules})
	}

	return config, nil
}

func createRule(rtype string, config RuleConfig) (TagRule, error) {
	var rule TagRule = nil

	switch rtype {
	case "file_exists":
		rule = &FileExistsRule{}
	case "file_contains":
		rule = &FileContainsRule{}
	case "in_path":
		rule = &InPathRule{}
	default:
		return nil, fmt.Errorf("Unknown rule type: %s", rtype)
	}

	err := rule.SetConfig(config)
	return rule, err
}
