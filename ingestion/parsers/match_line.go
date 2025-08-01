package parsers

import (
	"encoding/json"
	"regexp"
)

// match parser parses the input line and returns parsed value if match found

func MatchGroup(input string, cfg ExtractorConfig) (any, error) {
	re := regexp.MustCompile(cfg.Pattern)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1], nil // return as string
	}
	return nil, nil
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
