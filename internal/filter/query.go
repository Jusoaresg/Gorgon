package filter

import (
	"fmt"
	"strings"
)

// ExpandQuery expands a search pattern into a literal query string for a
// single name. Unlike Compile, the output is plain text used to build
// Prowlarr queries, so there is no regex escaping and {alias} is replaced by
// exactly the given name.
func ExpandQuery(pattern string, ctx Context, name string) (string, error) {
	var sb strings.Builder

	index := 0
	for _, match := range placeholderRe.FindAllStringSubmatchIndex(pattern, -1) {
		start, end := match[0], match[1]
		if start > index {
			sb.WriteString(pattern[index:start])
		}

		ph := pattern[match[2]:match[3]]

		width := 0
		if match[4] != -1 {
			width = len(pattern[match[4]:match[5]])
		}

		switch ph {
		case "alias":
			if name == "" {
				return "", fmt.Errorf("{alias} used but no name provided")
			}
			sb.WriteString(name)
		case "show":
			sb.WriteString(ctx.Show)
		case "season":
			sb.WriteString(formatNumber(ctx.Season, width))
		case "episode", "ep":
			sb.WriteString(formatNumber(ctx.Episode, width))
		case "absolute":
			return "", fmt.Errorf("placeholder {absolute} is reserved but not implemented yet")
		default:
			return "", fmt.Errorf("unknown placeholder {%s}", ph)
		}

		index = end
	}

	if index < len(pattern) {
		sb.WriteString(pattern[index:])
	}

	return sb.String(), nil
}
