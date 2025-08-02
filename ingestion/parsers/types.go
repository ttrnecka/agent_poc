package parsers

type ExtractFunc func(input string, cfg ExtractorConfig) (any, error)

type ExtractorConfig struct {
	Method  string `yaml:"method"`            // e.g., "match_group", "parse_json"
	Pattern string `yaml:"pattern,omitempty"` // for MatchGroup
	Path    string `yaml:"path,omitempty"`    // for JSONPath, etc.
}

type MapperConfig struct {
	Pattern string `yaml:"pattern"` // for endpoint matching
}

type Config struct {
	// first map points to endpoints
	// 2nd map points to object keys as they should be saved in object
	// example:
	// extractors:
	// version:
	//   fabric_os:
	//     method: match_group
	//     pattern: "(?m)^Fabric OS:\\s*(.+)$"
	//
	// defines version endpoint that should be parsed into fabric_os key in the final object
	// the version endpoint should be parsed using match_group extractor function using the pattern
	Extractors map[string]map[string]ExtractorConfig `yaml:"extractors"`
	// map again points to extractor endpoint to use if the edpoint matches the pattern
	Mappers map[string]MapperConfig `yaml:"mappers"`
}
