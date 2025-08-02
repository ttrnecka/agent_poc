package parsers

import (
	"encoding/json"
	"errors"
	"regexp"
)

var Extractors = map[string]ExtractFunc{
	"match_group":     MatchGroup,
	"match_group_all": MatchNamedGroupsAll,
	"parse_json":      ParseJSON,
}

// match parser parses the input line and returns parsed value if match found

func MatchGroup(input string, cfg ExtractorConfig) (any, error) {
	if cfg.Pattern == "" {
		return nil, errors.New("missing regex pattern")
	}

	re := regexp.MustCompile(cfg.Pattern)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1], nil // return as string
	}
	return nil, nil
}

func MatchNamedGroupsAll(input string, cfg ExtractorConfig) (any, error) {
	if cfg.Pattern == "" {
		return nil, errors.New("missing regex pattern")
	}

	re := regexp.MustCompile(cfg.Pattern)
	names := re.SubexpNames()
	matches := re.FindAllStringSubmatch(input, -1)

	var result []map[string]any

	for _, match := range matches {
		entry := make(map[string]any)
		for i, val := range match {
			if i == 0 {
				continue // skip the full match
			}
			if names[i] != "" {
				entry[names[i]] = val
			}
		}
		if len(entry) > 0 {
			result = append(result, entry)
		}
	}

	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func ParseJSON(input string, cfg ExtractorConfig) (any, error) {
	var data any
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return nil, err
	}

	// Example: simple traversal for fixed path
	if m, ok := data.(map[string]any); ok {
		if info, ok := m["info"].(map[string]any); ok {
			return info["meta"], nil // may be string, number, etc.
		}
	}

	return nil, nil
}
