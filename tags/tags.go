package tags

import (
	"fmt"

	"github.com/mecha/tags/log"
	"github.com/mecha/tags/rules"
)

type (
	Tag struct {
		Rules []rules.Rule
	}
)

func IsMatch(dir string, tag *Tag) bool {
	match, err := tag.CheckDir(dir)

	if err != nil {
		log.Error(err.Error())
		return false
	}

	return match
}

func (tag *Tag) CheckDir(dir string) (bool, error) {
	for _, rule := range tag.Rules {
		match, err := rule.Evaluate(dir)

		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}

	return false, nil
}

func (tag *Tag) findRule(ruleType string) int {
	for i, rule := range tag.Rules {
		if rules.GetType(rule) == ruleType {
			return i
		}
	}

	return -1
}

func (tag *Tag) AddRule(ruleType string, args []string) error {
	idx := tag.findRule(ruleType)

	if idx >= 0 {
		tag.Rules[idx].Add(args)
	} else {
		newRule, err := rules.New(ruleType)

		if err != nil {
			return err
		}

		newRule.Add(args)
		tag.Rules = append(tag.Rules, newRule)
	}

	return nil
}

func (tag *Tag) DelRule(ruleType string, args []string) error {
	idx := tag.findRule(ruleType)

	if idx >= 0 {
		tag.Rules[idx].Del(args)
		return nil
	}

	return fmt.Errorf("Rule \"%s\" not found.", ruleType)
}
