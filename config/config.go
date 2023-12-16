package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/mecha/tags/rules"
	"github.com/mecha/tags/tags"
	"github.com/mecha/tags/utils"
)

type (
	Config     map[string]TagConfig
	TagConfig  map[string]RuleConfig
	RuleConfig map[string]interface{}
)

func DefaultPath() string {
	defPath := xdg.ConfigHome() + "/tags/rules.json"

	envPath, ok := os.LookupEnv("TAGS_CONFIG")
	if !ok {
		return defPath
	}

	expPath, err := utils.ExpandTilde(envPath)
	if err != nil {
		return defPath
	}

	return expPath
}

func Read(path string) (map[string]tags.Tag, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	cfg, err := parseFile(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]tags.Tag, 0)

	for name, tagCfg := range cfg {
		rules, err := ruleListFromConfig(tagCfg)

		if err != nil {
			return nil, err
		}

		result[name] = tags.Tag{Rules: rules}
	}

	return result, nil
}

func Write(path string, data map[string]tags.Tag) error {
	config := Config{}

	for name, tag := range data {
		if len(tag.Rules) == 0 {
			continue
		}

		config[name] = make(TagConfig, 0)

		for _, rule := range tag.Rules {
			ruleType := ""
			switch rule.(type) {
			case *rules.FileExists:
				ruleType = "file_exists"
			case *rules.FileContains:
				ruleType = "file_contains"
			case *rules.InPath:
				ruleType = "in_path"
			default:
				return fmt.Errorf("Unknown rule type: %T", rule)
			}

			ruleCfg := rule.GetConfig()

			if len(ruleCfg) > 0 {
				config[name][ruleType] = ruleCfg
			}
		}
	}

	str, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(DefaultPath(), str, 0644)

	return nil
}

func ruleListFromConfig(cfg TagConfig) ([]rules.Rule, error) {
	ruleList := make([]rules.Rule, 0)

	for rType, rCfg := range cfg {
		rule, err := ruleFromConfig(rType, rCfg)
		if err != nil {
			return ruleList, err
		}

		ruleList = append(ruleList, rule)
	}

	return ruleList, nil
}

func ruleFromConfig(rType string, cfg RuleConfig) (rules.Rule, error) {
	var (
		rule rules.Rule = nil
		err  error      = nil
	)

	switch rType {
	case "file_exists":
		rule = &rules.FileExists{}
	case "file_contains":
		rule = &rules.FileContains{}
	case "in_path":
		rule = &rules.InPath{}
	default:
		return nil, fmt.Errorf("Unknown rule type: %s", rType)
	}

	err = rule.Load(cfg)

	if err != nil {
		return nil, err
	} else {
		return rule, nil
	}
}

func parseFile(path string) (Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result Config
	err = json.Unmarshal(raw, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
