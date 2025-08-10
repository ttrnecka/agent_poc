package api

import (
	"strconv"
	"strings"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// func output(file string) string {
// 	cmd_file := fmt.Sprintf("data/api/%s", file)
// 	if _, err := os.Stat(cmd_file); os.IsNotExist(err) {
// 		return fmt.Sprintf("error: file %s does not exist", cmd_file)
// 	}

// 	b, err := os.ReadFile(cmd_file) // just pass the file name
// 	if err != nil {
// 		logger.Error().Err(err).Msg("")
// 	}
// 	return string(b)
// }

// match reports whether path matches the given pattern, which is a
// path with '+' wildcards wherever you want to use a parameter. Path
// parameters are assigned to the pointers in vars (len(vars) must be
// the number of wildcards), which must be of type *string or *int.
func match(path, pattern string, vars ...interface{}) bool {
	for ; pattern != "" && path != ""; pattern = pattern[1:] {
		switch pattern[0] {
		case '+':
			// '+' matches till next slash in path
			slash := strings.IndexByte(path, '/')
			if slash < 0 {
				slash = len(path)
			}
			segment := path[:slash]
			path = path[slash:]
			switch p := vars[0].(type) {
			case *string:
				*p = segment
			case *int:
				n, err := strconv.Atoi(segment)
				if err != nil || n < 0 {
					return false
				}
				*p = n
			default:
				panic("vars must be *string or *int")
			}
			vars = vars[1:]
		case path[0]:
			// non-'+' pattern byte must match path byte
			path = path[1:]
		default:
			return false
		}
	}
	return path == "" && pattern == ""
}
