package parsers

import (
	"reflect"
	"testing"
)

func TestMatchGroup(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		config   ExtractorConfig
		expected any
	}{
		{
			name:  "simple string match",
			input: "Name: Alice",
			config: ExtractorConfig{
				Method:  "match_group",
				Pattern: `(?m)^Name:\s*(.+)$`,
			},
			expected: "Alice",
		},
		{
			name:  "number match",
			input: "Age: 42",
			config: ExtractorConfig{
				Method:  "match_group",
				Pattern: `(?m)^Age:\s*(\d+)$`,
			},
			expected: "42",
		},
		{
			name:  "no match",
			input: "No match here",
			config: ExtractorConfig{
				Method:  "match_group",
				Pattern: `(?m)^Missing:\s*(.+)$`,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, ok := Extractors[tt.config.Method]
			if !ok {
				t.Fatalf("unknown method: %s", tt.config.Method)
			}

			result, err := fn(tt.input, tt.config)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected: %#v, got: %#v", tt.expected, result)
			}
		})
	}
}

// func TestParseJSON(t *testing.T) {
// 	input := `{"info": {"meta": {"score": 88, "tags": ["x", "y"]}}}`
// 	expected := map[string]any{
// 		"score": float64(88),
// 		"tags":  []any{"x", "y"},
// 	}

// 	result, err := ParseJSON(input, ExtractorConfig{Method: "parse_json"})
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("expected: %#v, got: %#v", expected, result)
// 	}
// }
