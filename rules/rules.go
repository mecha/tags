package rules

import (
	"fmt"
)

type (
	Rule interface {
		fmt.Stringer
		Evaluate(dir string) (bool, error)
		Add(args []string) error
		Del(args []string) error
		Load(cfg map[string]interface{}) error
		GetConfig() map[string]interface{}
	}
)

func New(rType string) (Rule, error) {
	switch rType {
	case "file_exists":
		return &FileExists{}, nil
	case "file_contains":
		return &FileContains{}, nil
	case "in_path":
		return &InPath{}, nil
	default:
		return nil, fmt.Errorf("Unknown rule type: %s", rType)
	}
}

func GetType(rule Rule) string {
	switch rule.(type) {
	case *FileExists:
		return "file_exists"
	case *FileContains:
		return "file_contains"
	case *InPath:
		return "in_path"
	default:
		return ""
	}
}

func Add(rules []Rule, rType string, args []string) ([]Rule, error) {
	for _, r := range rules {
		if GetType(r) == rType {
			err := r.Add(args)
			if err != nil {
				return rules, err
			} else {
				break
			}
		}
	}

	return rules, nil
}

func Del(rules []Rule, rType string, args []string) ([]Rule, error) {
	if len(args) == 0 {
		return make([]Rule, 0), nil
	}

	for _, r := range rules {
		if GetType(r) == rType {
			err := r.Del(args)
			if err != nil {
				return rules, err
			} else {
				break
			}
		}
	}

	return rules, nil
}
